import {ApphostConn} from "../apphost/adapter.js";
import {bind} from "./bind";
import {call} from "./call";
import {hasParams} from "./query";
import {formatQueryParams} from "./params";

/**
 * Adds RPC implementation to {@link ApphostConn}. Compatible with astral query.
 *
 * @extends {ApphostConn}
 */
export class RpcConn extends ApphostConn {

  constructor(data) {
    super(data);
  }

  /**
   * Returns a copy of the object with each route bound as a method.
   *
   * @param {...string|object} routes
   * @return {RpcConn&object}
   *
   * @example
   * const conn = rpc.target(id).conn()
   * const api = conn.bind("foo", "bar")
   * await api.foo()
   * await api.bar("baz", 1, true)
   */
  bind(...routes) {
    return bind(this, routes);
  }

  /**@override*/
  copy() {
    return this;
  }

  /**@override*/
  call(port, ...params) {
    const c = call(this, this.#sub(port), ...params);
    c.inner.single = false
    return c
  }

  #sub(port) {
    if (hasParams(this.query)) throw `cannot nest connection for complete query ${chunks}`
    return port
  }

  /**
   * Adds response mapper to this connection.
   *
   * @param {(any) => any} f - mapping function
   * @returns RpcConn - new instance
   */
  map(f) {
    if (this.mapper) {
      const map = this.mapper
      this.mapper = arg => f(map(arg))
      return this
    }
    this.mapper = f
    return this
  }

  /**
   * Writes given method with params to this connection.
   *
   * @async
   * @param {string} method
   * @param {...any} params
   * @returns {Promise<RpcConn>}
   */
  async conn(method, ...params) {
    let cmd = method ? method : ""
    if (params.length > 0) {
      if (cmd) cmd += '?'
      cmd += formatQueryParams(params)
    }
    if (cmd) await this.writeLn(cmd)
    return this
  }

  /**
   * Encodes data into JSON string and writes as line.
   *
   * @async
   * @param {any} data
   * @returns {Promise<undefined>}
   * @throws {string} - IO error message
   */
  async encode(data) {
    let json = JSON.stringify(data)
    if (json === undefined) json = '{}'
    return await super.writeLn(json)
  }

  /**
   * Reads line and parses it into object.
   *
   * @async
   * @returns {Promise<object>} - parsed JSON object
   * @throws {string} - parsed error message
   */
  async decode() {
    const resp = await this.readLn()
    const parsed = JSON.parse(resp)
    if (parsed === null) return null
    if (parsed.error) throw parsed.error
    return parsed
  }

  /**
   * Returns first successfully decoded value or null.
   *
   * @async
   * @returns {Promise<any|null>}
   */
  async request() {
    const map = this.mapper
    this.result = null
    for (; ;) {
      const next = await this.decode()
      if (next === undefined) continue
      if (next === null) return this.result
      this.result = next
      if (!map) return next
      const last = await map(next)
      if (last === undefined) continue
      if (last === null) return this.result
      return last
    }
  }

  /**
   * Collects all decoded values until first null occurrence.
   *
   * @returns {Promise<any[]>} - collected values
   */
  async collect() {
    const map = this.mapper ? this.mapper : null
    let push
    if (!map) push = next => this.result.push(next)
    else push = async (next) => {
      next = await map.call(this, next)
      if (next === null) return this.result
      if (next) this.result.push(next)
    }
    this.result = []
    for (; ;) {
      let next = await this.decode()
      if (next === null) return this.result
      push(next)
    }
  }
}
