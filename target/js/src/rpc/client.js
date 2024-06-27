import {ApphostClient} from "../apphost/adapter.js";
import {RpcConn} from "./conn.js";
import {serve} from "./serve.js";
import {call} from "./call";
import {bind} from "./bind";

export class RpcClient extends ApphostClient {

  bind(...routes) {
    return bind(this, routes);
  }

  copy(data) {
    return Object.assign(new RpcClient(), {...this, ...data});
  }

  target(id) {
    this.targetId = id
    return this
  }

  call(port, ...params) {
    return call(this, port, ...params);
  }

  async conn(port, ...params) {
    const query = formatQuery(port, params)
    const conn = await super.query(query, this.targetId)
    return new RpcConn(conn)
  }

  async serve(ctx) {
    return await serve(this, ctx);
  }
}

function formatQuery(port, params) {
  let query = port
  if (params.length > 0) {
    query += '?' + JSON.stringify(params)
  }
  return query
}

