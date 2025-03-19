import {ApphostClient} from "../apphost/adapter.js";
import {RpcConn} from "./conn.js";
import {serve} from "./serve.js";
import {call} from "./call";
import {bind} from "./bind";
import {formatQueryParams} from "./params";

export class RpcClient extends ApphostClient {

  bind(...routes) {
    return bind(this, routes);
  }

  copy(data) {
    return Object.assign(new RpcClient(), {...this, ...data});
  }

  target(id) {
    return this.copy({targetId: id})
  }

  call(port, ...params) {
    port = port ? port : ""
    return call(this, port, ...params);
  }

  async conn(port, ...params) {
    port = port ? port : ""
    const query = formatQuery(port, params)
    const conn = await super.query(this.targetId, query)
    return new RpcConn(conn)
  }

  async serve(ctx) {
    return await serve(this, ctx);
  }
}

function formatQuery(port, params) {
  let query = port
  if (params.length > 0) {
    query += '?' + formatQueryParams(params)
  }
  return query
}
