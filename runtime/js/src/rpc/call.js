import bindings from "../bindings";

/**
 * @param {RpcClient | RpcConn} client
 * @param {string} port
 * @param {any[]} params
 */
export function call(client, port, ...params) {
  const call = new RpcCall(client, port, params)
  let f = async (...params) => {
    return await call.request(...params);
  }
  return Object.assign(f, {
    inner: call,
    map: (...args) => call.map(...args),
    filter: (...args) => call.filter(...args),
    request: async (...args) => await call.request(...args),
    collect: async (...args) => await call.collect(...args),
    conn: async (...args) => await call.conn(...args),
  })
}

class RpcCall {

  mapper = arg => arg
  params = []
  single = true

  /**
   * @param {RpcClient | RpcConn} client
   * @param {string} port
   * @param {any[]} params
   */
  constructor(client, port, params) {
    this.client = client
    this.port = port
    this.params = Array.isArray(params) ? params : params ? [params] : []
  }

  map(f) {
    const map = this.mapper
    this.mapper = arg => f(map(arg))
    return this
  }

  filter(f) {
    return this.map(arg => {
      if (f(arg)) return arg
    })
  }

  async request(...params) {
    if (params.length > 0) this.params = params
    return await this.#consume(async conn => await conn.request(...params));
  }

  async collect(...params) {
    if (params.length > 0) this.params = params
    return await this.#consume(async conn => await conn.collect(...params));
  }

  async #consume(f) {
    const conn = await this.conn()
    conn.mapper = this.mapper
    this.result = await f(conn)
    this.mapper = a => a // reset mapper between requests
    if (this.single) await conn.close().catch(bindings.log)
    return this.result
  }

  async conn(...params) {
    const args = params.length > 0 ? params : this.params
    return this.client.conn(this.port, ...args);
  }
}
