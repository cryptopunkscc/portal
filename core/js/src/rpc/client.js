import {ApphostClient} from "../apphost/adapter.js";
import {RpcConn} from "./conn.js";
import {serve} from "./serve.js";
import {call} from "./call";
import {bind} from "./bind";
import {formatQuery} from "./query";

/**
 * Adds RPC implementation to {@link ApphostClient}. Compatible with astral query.
 *
 * @extends {ApphostClient}
 */
export class RpcClient extends ApphostClient {

  /**
   * Returns a copy with assigned data.
   *
   * @param data {object}
   * @return RpcClient
   */
  copy(data) {
    return Object.assign(new RpcClient(), {...this, ...data});
  }

  /**
   * Returns a copy with given target id.
   *
   * @param {string} id or alias
   * @return RpcClient
   */
  target(id) {
    return this.copy({targetId: id})
  }

  /**
   * Returns a copy with each route bound as a method.
   *
   * @param {...(string|object)} routes
   * @return {RpcClient&object}
   * @example
   * // A copy of rpc object with target id where foo method is bound to the "foo" path and fiz method is bound to the "bar.baz.fiz" path.
   * const api = rpc.target(id).bind("foo", {"bar.baz": ["fiz"]})
   * await api.foo()
   * await api.fiz("yolo", 1, true)
   */
  bind(...routes) {
    return bind(this, routes);
  }

  /**
   * Creates a call object from this {@link RpcClient}.
   *
   * @param {string} port
   * @param {...arg} params
   * @returns {function&RpcCall}
   */
  call(port, ...params) {
    port = port ? port : ""
    return call(this, port, ...params);
  }

  /**
   * Opens new {@link RpcConn}.
   *
   * @async
   * @param {string} port
   * @param {...arg} params
   * @returns {RpcConn}
   */
  async conn(port, ...params) {
    port = port ? port : ""
    const query = formatQuery(port, params)
    const conn = await super.query(this.targetId, query)
    return new RpcConn(conn)
  }

  /**
   * Serves given API via apphost.
   *
   * @async
   * @param {object} api
   * @returns {Promise<void>}
   */
  async serve(api) {
    return await serve(this, api);
  }
}
