import {ApphostClient} from "../apphost/adapter.js";
import {RpcConn} from "./conn.js";
import {serve} from "./serve.js";

export class RpcClient extends ApphostClient {

  constructor(targetId, methods) {
    super();
    this.targetId = targetId
    this.boundMethods = methods
  }

  async query(query){
    let conn = await super.query(query, this.targetId)
    conn = new RpcConn(conn)
    return conn
  }

  async serve(ctx) {
    await serve(this, ctx)
  }

  async call(query, ...params) {
    if (params) {
      query += '?' + JSON.stringify(params)
    }
    const conn = await super.query(query, this.targetId)
    return new RpcConn(conn)
  }

  async request(query, ...params) {
    const conn = await this.call(query, ...params)
    const response = await conn.decode()
    conn.close().catch(log)
    return response
  }

  caller(query) {
    return async (...params) => await this.call(query, ...params)
  }

  requester(query) {
    return async (...params) => await this.request(query, ...params)
  }

  target(id) {
    return new RpcClient(id)
  }

  bind(...methods) {
    const copy = new RpcClient(this.targetId, methods)
    for (let method of methods) {
      this[method] = copy.requester(method)
    }
    return copy
  }
}
