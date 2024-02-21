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
    const json = await bindings.astral_query(node, query)
    const data = JSON.parse(json)
    return new AppHostConn(data, query)
  }

  async queryName(node, query) {
    const json = await bindings.astral_query_name(node, query)
    const data = JSON.parse(json)
    return new AppHostConn(data, query)
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
    const json = await bindings.astral_conn_accept(this.port)
    const data = JSON.parse(json)
    return new AppHostConn(data)
  }

  async close() {
    await bindings.astral_service_close(this.port)
  }
}

class AppHostConn {
  constructor(data) {
    this.id = data.id
    this.query = data.query
  }

  async read() {
    return await bindings.astral_conn_read(this.id)
  }

  async write(data) {
    return await bindings.astral_conn_write(this.id, data)
  }

  async close() {
    await bindings.astral_conn_close(this.id)
  }
}
