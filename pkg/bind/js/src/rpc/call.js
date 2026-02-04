import bindings from "../bindings";

/**
 * Combines RPC client interface with its request function.
 *
 * @param {RpcClient | RpcConn} client
 * @param {string} port
 * @param {...any} params
 * @returns {function & RpcCall}
 */
export function call(client, port, ...params) {
  const call = new RpcCall(client, port, params)
  let f = call.request.bind(call)
  return Object.assign(f, {
    inner: call,
    map: (f) => call.map(f),
    filter: (f) => call.filter(f),
    request: async (...args) => await call.request(...args),
    collect: async (...args) => await call.collect(...args),
    conn: async (...args) => await call.conn(...args),
  })
}

/**
 * Represents a further intention for opening connections on a specific port with a given params.
 */
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

  /**
   * Adds response mapper.
   *
   * @param {(any) => any} f - mapping function
   * @returns RpcCall - new instance
   */
  map(f) {
    const map = this.mapper
    this.mapper = arg => f(map(arg))
    return this
  }

  /**
   * Adds response filter.
   *
   * @param {(any) => boolean} f - filtering function
   * @returns RpcCall - new instance
   */
  filter(f) {
    return this.map(arg => {
      if (f(arg)) return arg
    })
  }

  /**
   * Creates new {@link RpcConn} and calls {@link RpcConn.request}.
   *
   * @async
   * @param {...any} params - connection params
   * @returns {Promise<any>}
   */
  async request(...params) {
    if (params.length > 0) this.params = params
    return await this.#consume(async conn => await conn.request());
  }

  /**
   * Creates new {@link RpcConn} and calls {@link RpcConn.collect}.
   *
   * @async
   * @param {...any} params - connection params
   * @returns {Promise<any[]>}
   */
  async collect(...params) {
    if (params.length > 0) this.params = params
    return await this.#consume(async conn => await conn.collect());
  }

  async #consume(f) {
    const conn = await this.conn()
    conn.mapper = this.mapper
    this.result = await f(conn)
    this.mapper = a => a // reset mapper between requests
    if (this.single) await conn.close().catch(bindings.log)
    return this.result
  }

  /**
   * Returns new {@link RpcConn}.
   *
   * @async
   * @param {...any} params - connection params
   * @returns {Promise<RpcConn>}
   */
  async conn(...params) {
    const args = params.length > 0 ? params : this.params
    return this.client.conn(this.port, ...args);
  }
}
