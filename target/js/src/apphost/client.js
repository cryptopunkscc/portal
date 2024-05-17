import {bindings} from "../bindings";

// ================== Object oriented adapter ==================

export class AppHostClient {
  async register(service) {
    await bindings.astral_service_register(service)
    return new AppHostListener(service)
  }

  async query(identity, query) {
    const json = await bindings.astral_query(identity, query)
    const data = JSON.parse(json)
    return new AppHostConn(data, query)
  }

  async queryName(name, query) {
    const json = await bindings.astral_query_name(name, query)
    const data = JSON.parse(json)
    return new AppHostConn(data, query)
  }

  async nodeInfo(id) {
    return await bindings.astral_node_info(id)
  }

  async resolve(name) {
    return await bindings.astral_resolve(name)
  }

  async interrupt() {
    await bindings.astral_interrupt()
  }
}

export class AppHostListener {
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

export class AppHostConn {
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
