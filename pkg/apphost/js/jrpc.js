// ================== RPC extensions ==================

// Bind RPC api of service associated to this connection
AppHostConn.prototype.bindRpc = async function () {
  await astral_rpc_conn_bind_api(this)
}

AppHostClient.prototype.bindRpc = async function (node, service) {
  await astral_rpc_client_bind_api(this, node, service)
  return this
}

async function astral_rpc_conn_bind_api(conn) {
  // request api methods
  const cmd = JSON.stringify(["api"])
  await conn.write(cmd)

  // read api methods
  const api = await conn.read()
  const methods = JSON.parse(api)
  log(conn.query + " " + conn.id + ": == " + api)

  // bind methods
  for (let method of methods) {
    conn[method] = async (...data) => {
      const cmd = JSON.stringify([method, ...data])
      log(conn.query + " " + conn.id + ": => " + cmd)
      await conn.write(cmd)
      const resp = await conn.read()
      const json = JSON.parse(resp)
      log(conn.query + " " + conn.id + ": <= " + JSON.stringify([method, json]))
      return json
    }
  }
}

async function astral_rpc_client_bind_api(client, node, service) {
  // request api methods
  const query = service + JSON.stringify(["api"])
  const conn = await client.query(node, query)

  // read api methods
  const api = await conn.read()
  const methods = JSON.parse(api)
  log(service + " " + conn.id + ": == " + api)
  conn.close()

  // bind methods
  for (let method of methods) {
    client[method] = async (...data) => {
      const cmd = JSON.stringify([method, ...data])
      const conn = await client.query(node, service + cmd)
      log(service + " " + conn.id + ": => " + cmd)
      const resp = await conn.read()
      const json = JSON.parse(resp)
      conn.close().catch()
      log(service + " " + conn.id + ": <= " + JSON.stringify([method, json]))
      return json
    }
  }
}

// Bind RPC service to given name
AppHostClient.prototype.bindRpcService = async function (service) {
  await astral_rpc_bind_srv.call(this, service)
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
  astral_rpc_listen.call(srv, listener).then(log)
}

async function astral_rpc_listen(listener) {
  for (; ;) {
    const conn = await listener.accept()
    log(conn.query + " " + conn.id + ": accepted")
    astral_rpc_handle.call(this, conn)
  }
}

async function astral_rpc_handle(conn) {
  try {
    let str = conn.query.slice(this.name.length)
    const single = str.length > 0
    for (; ;) {
      if (!single) {
        str = await conn.read()
      }
      log(this.name + " " + conn.id + ": " + str)
      const query = JSON.parse(str)
      const method = query[0]
      const args = query.slice(1)
      const result = await this[method](...args)
      await conn.write(JSON.stringify(result))
      if (single) {
        break
      }
    }
  } catch (e) {
    log(conn.query + " " + conn.id + ": " + e)
  }
}
