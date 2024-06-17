import {ApphostClient} from "../apphost/adapter.js";
import {RpcConn} from "./conn.js";
import {serve} from "./serve.js";
import {bindings} from "../bindings"

export class RpcClient extends ApphostClient {

  constructor() {
    super();
    this.port = ""
  }

  async serve(ctx) {
    await serve(this, ctx)
  }

  async call(query, ...params) {
    if (params.length > 0) {
      query += '?' + JSON.stringify(params)
    }
    const conn = await super.query(query, this.targetId)
    return new RpcConn(conn)
  }

  async request(query, ...params) {
    const conn = await this.call(query, ...params)
    const response = await conn.decode()
    conn.close().catch(bindings.log)
    return response
  }

  caller(query) {
    return async (...params) => await this.call(query, ...params)
  }

  requester(query) {
    return async (...params) => await this.request(query, ...params)
  }

  observer(query) {
    return async (...params) => {
      if (typeof params[params.length - 1] !== "function") {
        return await this.request(query, ...params)
      }
      const consume = params.pop()
      const conn = await this.call(query, ...params)
      let last
      try {
        last = await conn.observe(consume)
      }
      finally {
        conn.close().finally()
      }
      bindings.log("observer last", last)
      return last
    }
  }

  copy(data) {
    return Object.assign(new RpcClient(), {...this, ...data})
  }

  target(id) {
    return this.copy({targetId: id})
  }

  bind(route, ...methods) {
    const copy = this.copy()
    for (let method of methods) {
      const collect = method[0] === '*'
      if (collect) {
        method = method.substring(1) // drop * prefix
      }
      if (this[method]) {
        throw `method '${method}' already exist`
      }
      const port = [route, method].join('.')
      if (collect) {
        copy[method] = copy.observer(port)
      } else {
        copy[method] = copy.requester(port)
      }
    }
    return copy
  }
}
