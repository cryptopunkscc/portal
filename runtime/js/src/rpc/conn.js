import {ApphostConn} from "../apphost/adapter.js";
import {bind} from "./bind";
import {call} from "./call";
import {hasParams} from "./query";
import {formatQueryParams} from "./params";

export class RpcConn extends ApphostConn {

  constructor(data) {
    super(data);
  }

  #sub(port) {
    if (hasParams(this.query)) throw `cannot nest connection for complete query ${chunks}`
    return port
  }

  bind(...routes) {
    return bind(this, routes);
  }

  copy() {
    return this;
  }

  call(port, ...params) {
    const c = call(this, this.#sub(port), ...params);
    c.inner.single = false
    return c
  }

  map(f) {
    if (this.mapper) {
      const map = this.mapper
      this.mapper = arg => f(map(arg))
    } else {
      this.mapper = f
    }
    return this
  }

  async conn(method, ...params) {
    let cmd = method ? method : ""
    if (params.length > 0) {
      if (cmd) cmd += '?'
      cmd += formatQueryParams(params)
    }
    if (cmd) await this.writeLn(cmd)
    return this
  }

  async encode(data) {
    let json = JSON.stringify(data)
    if (json === undefined) json = '{}'
    return await super.writeLn(json)
  }

  async decode() {
    const resp = await this.readLn()
    const parsed = JSON.parse(resp)
    if (parsed === null) return null
    if (parsed.error) throw parsed.error
    return parsed
  }

  async request(...params) {
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
   * Collects all decoded values mapped as not null until decodes null or maps into undefined.
   *
   * @param {...any} params
   * @returns {Promise<any[]>}
   */
  async collect(...params) {
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
