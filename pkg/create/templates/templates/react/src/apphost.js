let builder = []

// ================== Wails bindings adapter ==================

/* eslint-disable */
const _wails_platform = () => typeof window['go'] === "undefined" ? undefined : "wails"

/* eslint-disable */
const _wails_bindings = () => ({
  astral_conn_accept: window['go']['main']['Adapter']['ConnAccept'],
  astral_conn_close: window['go']['main']['Adapter']['ConnClose'],
  astral_conn_read: window['go']['main']['Adapter']['ConnRead'],
  astral_conn_write: window['go']['main']['Adapter']['ConnWrite'],
  astral_node_info: window['go']['main']['Adapter']['NodeInfo'],
  astral_query: window['go']['main']['Adapter']['Query'],
  astral_query_name: window['go']['main']['Adapter']['QueryName'],
  astral_resolve: window['go']['main']['Adapter']['Resolve'],
  astral_service_close: window['go']['main']['Adapter']['ServiceClose'],
  astral_service_register: window['go']['main']['Adapter']['ServiceRegister'],
  sleep: window['go']['main']['Adapter']['Sleep'],
  log: window['go']['main']['Adapter']['LogArr'],
})

builder.push({
  platform: _wails_platform(),
  bindings: _wails_bindings,
})

// ================== Android bindings adapter ==================

/* eslint-disable */
const _android_platform = () => typeof _app_host === "undefined" ? undefined : "android"

/* eslint-disable */
const _android_bindings = () => {

  const _awaiting = new Map()

  window._resolve = (id, value) => {
    _awaiting.get(id)[0](value)
    _awaiting.delete(id)
  }

  window._reject = (id, error) => {
    _awaiting.get(id)[1](error)
    _awaiting.delete(id)
  }

  const _promise = (block) =>
    new Promise((resolve, reject) =>
      _awaiting.set(block(), [resolve, reject]))

  return {
    astral_node_info: (arg1) => _promise(() => _app_host.nodeInfo(arg1)).then(v => JSON.parse(v)),
    astral_conn_accept: (arg1) => _promise(() => _app_host.connAccept(arg1)),
    astral_conn_close: (arg1) => _promise(() => _app_host.connClose(arg1)),
    astral_conn_read: (arg1) => _promise(() => _app_host.connRead(arg1)),
    astral_conn_write: (arg1, arg2) => _promise(() => _app_host.connWrite(arg1, arg2)),
    astral_query: (arg1, arg2) => _promise(() => _app_host.query(arg1, arg2)),
    astral_query_name: (arg1, arg2) => _promise(() => _app_host.queryName(arg1, arg2)),
    astral_resolve: (arg1) => _promise(() => _app_host.resolve(arg1)),
    astral_service_close: (arg1) => _promise(() => _app_host.serviceClose(arg1)),
    astral_service_register: (arg1) => _promise(() => _app_host.serviceRegister(arg1)),
    sleep: (arg1) => _promise(() => _app_host.sleep(arg1)),
    log: (arg1) => _app_host.logArr(JSON.stringify(arg1)),
  }
}

builder.push({
  platform: _android_platform(),
  bindings: _android_bindings,
})

/* eslint-disable */
const _default_bindings = () => ({
  astral_conn_accept: _astral_conn_accept,
  astral_conn_close: _astral_conn_close,
  astral_conn_read: _astral_conn_read,
  astral_conn_write: _astral_conn_write,
  astral_node_info: _astral_node_info,
  astral_query: _astral_query,
  astral_query_name: _astral_query_name,
  astral_resolve: _astral_resolve,
  astral_service_close: _astral_service_close,
  astral_service_register: _astral_service_register,
  sleep: _sleep,
  log: _log,
})

builder.push({
  platform: "default",
  bindings: _default_bindings,
})

const platform = function () {
  for (let next of builder) {
    if (next.platform) {
      return next.platform
    }
  }
}()

const bindings = function () {
  for (let next of builder) {
    if (next.platform) {
      return next.bindings()
    }
  }
}()

// ================== Static functions adapter ==================

const log = (...arg1) => bindings.log(arg1)
const sleep = (arg1) => bindings.sleep(arg1)

// ================== Object oriented adapter ==================

class AppHostClient {
  async register(service) {
    await bindings.astral_service_register(service)
    return new AppHostListener(service)
  }

  async query(node, query) {
    const conn = await bindings.astral_query(node, query)
    return new AppHostConn(conn, query)
  }

  async queryName(node, query) {
    const conn = await bindings.astral_query_name(node, query)
    return new AppHostConn(conn, query)
  }

  async nodeInfo(id) {
    return await bindings.astral_node_info(id)
  }

  async resolve(name) {
    return await bindings.astral_resolve(name)
  }
}

class AppHostListener {
  constructor(port) {
    this.port = port
  }

  async accept() {
    const conn = await bindings.astral_conn_accept(this.port)
    return new AppHostConn(conn, this.port)
  }

  async close() {
    await bindings.astral_service_close(this.port)
  }
}

class AppHostConn {
  constructor(conn, port) {
    this.conn = conn
    this.port = port
  }

  async read() {
    return await bindings.astral_conn_read(this.conn)
  }

  async write(data) {
    return await bindings.astral_conn_write(this.conn, data)
  }

  async close() {
    await bindings.astral_conn_close(this.conn)
  }
}

const appHost = new AppHostClient()

// ================== RPC extensions ==================

// Bind RPC api of service associated to this connection
AppHostConn.prototype.bindRpc = async function () {
  await astral_rpc_bind_api(this)
}

async function astral_rpc_bind_api(conn) {
  await conn.write(JSON.stringify(["api"]))
  const api = await conn.read()
  log(conn.port + " " + conn.conn + ": == " + api)
  const methods = JSON.parse(api)
  for (let method of methods) {
    conn[method] = async (...data) => {
      const cmd = JSON.stringify([method, ...data])
      log(conn.port + " " + conn.conn + ": => " + cmd)
      await conn.write(cmd)
      const resp = await conn.read()
      const json = JSON.parse(resp)
      log(conn.port + " " + conn.conn + ": <= " + JSON.stringify([method, json]))
      return json
    }
  }
}

// Bind RPC service to given name
AppHostClient.prototype.bindRpc = async function (service, name) {
  await astral_rpc_bind_srv.call(this, service, name)
}

async function astral_rpc_bind_srv(Service, name) {
  const props = Object.getOwnPropertyNames(Service.prototype)
  if (props[0] !== "constructor") throw new Error("Service must have a constructor")
  const methods = props.slice(1, props.length)
  methods.push("api")
  Service.prototype.api = async () => {
    return methods
  }
  const listener = await this.register(name)
  log("listen " + name)
  astral_rpc_listen.call(new Service(), listener).then(log)
}

async function astral_rpc_listen(listener) {
  for (; ;) {
    const conn = await listener.accept()
    log(conn.port + " " + conn.conn + ": accepted")
    astral_rpc_handle.call(this, conn)
  }
}

async function astral_rpc_handle(conn) {
  try {
    for (; ;) {
      const str = await conn.read()
      log(conn.port + " " + conn.conn + ": " + str)
      const query = JSON.parse(str)
      const method = query[0]
      const args = query.slice(1)
      const result = await this[method](...args)
      await conn.write(JSON.stringify(result))
    }
  } catch (e) {
    log(conn.port + " " + conn.conn + ": " + e)
  }
}

// ================== Exports ==================

export default appHost

export {
  platform,
  log,
  sleep,
}

