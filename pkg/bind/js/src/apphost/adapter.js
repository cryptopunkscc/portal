import {bindings} from "../bindings";

// ================== Object oriented adapter ==================

export class ApphostClient {
  async register() {
    await bindings.astral_service_register()
    return new ApphostListener()
  }

  async query(target, query) {
    const data = await bindings.astral_query(target, query);
    return new ApphostConn(data, query)
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

export class ApphostListener {
  constructor() {
  }

  async accept() {
    const data = await bindings.astral_conn_accept()
    return new ApphostConn(data)
  }

  async close() {
    await bindings.astral_service_close()
  }
}

export class ApphostConn {
  constructor(data) {
    this.id = data.id
    this.query = data.query
    this.remoteId = data.remoteId
  }

  async read(buffer) {
    try {
      return await bindings.astral_conn_read(this.id, buffer)
    } catch (e) {
      this.done = true
      if (e === "EOF") {
        return -1
      }
      throw e
    }
  }

  async write(data) {
    try {
      return await bindings.astral_conn_write(this.id, data)
    } catch (e) {
      this.done = true
      throw e
    }
  }

  async readLn() {
    try {
      return await bindings.astral_conn_read_ln(this.id)
    } catch (e) {
      this.done = true
      throw e
    }
  }

  async writeLn(data) {
    try {
      return await bindings.astral_conn_write_ln(this.id, data)
    } catch (e) {
      this.done = true
      throw e
    }
  }

  async close() {
    if (!this.done) {
      this.done = true
      await bindings.astral_conn_close(this.id)
    }
  }
}
