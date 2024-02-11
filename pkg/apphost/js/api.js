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
