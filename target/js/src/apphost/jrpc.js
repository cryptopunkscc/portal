import {AppHostClient, AppHostConn} from "./client";
import {bindings} from "../bindings.js";

const {log} = bindings

export * from "./client";

// ================== RPC extensions ==================

AppHostConn.prototype.jrpcCall = async function (method, ...data) {
  let cmd = method
  if (data.length > 0) {
    cmd += "?" + JSON.stringify(data)
  }
  log(this.id + " conn => " + this.query + "." + cmd)
  await this.write(cmd + '\n')
}

AppHostConn.prototype.readJson = async function (method) {
  const resp = await this.read()
  const json = JSON.parse(resp)
  if (method !== undefined) {
    log(this.id + " conn <= " + this.query  + ":" + resp.trimEnd())
  }
  return json
}

AppHostConn.prototype.writeJson = async function (data) {
  const json = JSON.stringify(data)
  log(this.id + " conn => " + this.query + ":" + json.trimEnd())
  await this.write(json + '\n')
}

AppHostConn.prototype.jsonReader = async function (method) {
  const read = async () => await this.readJson(method)
  read.cancel = async () => await this.close()
  return read
}

// Bind RPC api of service associated to this connection
AppHostConn.prototype.bindRpc = async function () {
  await astral_rpc_conn_bind_api(this)
}

async function astral_rpc_conn_bind_api(conn) {
  // request api methods
  await conn.jrpcCall("api")

  // read api methods
  const methods = await conn.readJson("api")

  // bind methods
  for (let method of methods) {
    conn[method] = async (...data) => {
      await conn.jrpcCall(method, ...data)
      return await conn.readJson(method)
    }
  }

  // bind subscribe
  conn.subscribe = async (method, ...data) => {
    await conn.jrpcCall(method, ...data)
    return conn.jsonReader(method)
  }
}

AppHostClient.prototype.jrpcCall = async function (identity, service, method, ...data) {
  let cmd = service
  if (method) {
    cmd += "." + method
  }
  if (data.length > 0) {
    cmd += "?" + JSON.stringify(data)
  }
  const conn = await this.query(identity, cmd)
  log(conn.id + " client => " + cmd)
  return conn
}

AppHostClient.prototype.bindRpc = async function (identity, service) {
  await astral_rpc_client_bind_api(this, identity, service)
  return this
}

AppHostClient.prototype.rpcQuery = function (identity, port) {
  const client = this
  return async function (...data) {
    const conn = await client.jrpcCall(identity, port, "", ...data)
    const json = await conn.readJson(port)
    conn.close().catch(log)
    return json
  }
}

async function astral_rpc_client_bind_api(client, identity, service) {
  // request api methods
  const conn = await client.jrpcCall(identity, service, "api")

  // read api methods
  const methods = await conn.readJson("api")
  conn.close().catch(log)

  // bind methods
  for (let method of methods) {
    client[method] = async (...data) => {
      const conn = await client.jrpcCall(identity, service, method, ...data)
      const json = await conn.readJson(method)
      conn.close().catch(log)
      return json
    }
  }

  // bind subscribe
  client.subscribe = async (method, ...data) => {
    const conn = await client.jrpcCall(identity, service, method, ...data)
    return await conn.jsonReader(method)
  }
}

// Bind RPC service to given name
AppHostClient.prototype.bindRpcService = async function (service) {
  return await astral_rpc_bind_srv.call(this, service)
}

async function astral_rpc_bind_srv(Service) {
  const props = Object.getOwnPropertyNames(Service.prototype)
  if (props[0] !== "constructor") throw new Error("Service must have a constructor")
  const methods = props.slice(1, props.length)
  methods.push("api")
  Service.prototype.api = async () => {
    return methods
  }
  const srv = new Service()
  const listener = await this.register(srv.name + "*")
  log("listen " + srv.name)
  astral_rpc_listen.call(srv, listener).catch(log)
  return listener
}

async function astral_rpc_listen(listener) {
  for (; ;) {
    const conn = await listener.accept()
    log(conn.id + " service <= " + conn.query)
    astral_rpc_handle.call(this, conn).catch(log)
  }
}

async function astral_rpc_handle(conn) {
  try {
    let query = conn.query.slice(this.name.length)
    let method = query, args = []
    const single = query !== ''
    const write = async (data) => await conn.writeJson(data)
    const read = async (method) => await conn.readJson(method)

    for (; ;) {
      if (!single) {
        query = await conn.read();
        log(conn.id + " service <== " + query)
      }
      [method, args] = parseQuery(query)

      let result = await this[method](...args, write, read)
      if (result !== undefined) {
        await conn.writeJson(result)
      }
      if (single) {
        conn.close().catch(log)
        break
      }
    }
  } catch (e) {
    log(conn.id + " service !! " + conn.query + ":" + e)
    conn.close().catch(log)
  }
}

function parseQuery(query) {
  if (query[0] === '.') {
    query = query.slice(1)
  }
  let [method, payload] = query.split('?', 2)
  let args = []
  if (payload) {
    args = JSON.parse(payload)
  }
  return [method, args]
}
