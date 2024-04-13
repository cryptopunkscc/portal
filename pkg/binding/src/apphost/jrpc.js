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
  log(this.query + " " + this.id + ": => " + cmd)
  await this.write(cmd)
}

AppHostConn.prototype.readJson = async function (method) {
  const resp = await this.read()
  const json = JSON.parse(resp)
  if (method !== undefined) {
    log(this.query + " " + this.id + ": <= " + method  + ":" + resp)
  }
  return json
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
  log(service + " " + conn.id + ": => " + cmd)
  return conn
}

AppHostClient.prototype.bindRpc = async function (identity, service) {
  await astral_rpc_client_bind_api(this, identity, service)
  return this
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
    log(conn.query + " " + conn.id + ": accepted")
    astral_rpc_handle.call(this, conn).catch(log)
  }
}

async function astral_rpc_handle(conn) {
  try {
    const send = async (result) => await conn.write(JSON.stringify(result))

    let query = conn.query.slice(this.name.length)
    let method = query, args = []
    const single = query !== ''

    for (; ;) {
      if (!single) {
        query = await conn.read();
      }
      [method, args] = parseQuery(query)

      log(this.name + " " + conn.id + ": " + query)
      let result = await this[method](...args, send)
      if (result !== undefined) {
        result = JSON.stringify(result)
        await conn.write(result)
      }
      if (single) {
        conn.close().catch(log)
        break
      }
    }
  } catch (e) {
    log(conn.query + " " + conn.id + ": " + e)
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
