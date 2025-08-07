import {bindings} from "./bindings.js";
import {ApphostClient} from "./apphost/adapter.js";
import {RpcClient} from "./rpc/client.js";

/**
 * Logs given data to the native console.
 */
export const log = async any => await bindings.log(typeof any == 'object' ? JSON.stringify(any) : any)


export const {
  /**
   * Platform name constant string.
   */
  platform,

  /**
   * Closes process with given code.
   *
   * @async
   * @param {int} code - Exit code
   */
  exit,

  /**
   * Delays execution for a given milliseconds.
   *
   * @async
   * @param {bigint} millis
   */
  sleep,

} = bindings

/**
 * {@link ApphostClient} singleton - Provides bindings the native apphost implementation.
 *
 * @example
 * // You can register query listener to obtain a connection and begin communication.
 * const listener = await apphost.register()
 * const conn = await listener.accept()
 * log({id: conn.remoteId, query: conn.query})
 * const msg = await conn.decode()
 * await conn.encode({echo: msg})
 * await conn.close()
 * await listener.close()
 *
 * // You can query target by its alias or id to get new connection and read a data.
 * const conn = await.apphost.query(targetId, "foo:bar=baz")
 * const data = conn.decode()
 * await log(data)
 */
export const apphost = new ApphostClient();

/**
 * {@link RpcClient} singleton - Provides RPC API to the {@link ApphostClient} compatible with apphost query.
 *
 * @example
 * // You can register API handlers and inject optional dependencies.
 * rpc.serve({
 *   // Everything inside inject will be passed to the invoked handler in the first argument under the '$' key with assigned connection context.
 *   inject: {
 *     dispatcher: cmd => {...},
 *     state: {...}
 *   },
 *   handlers: {
 *     // Simple handler.
 *     func0: () => 0,
 *
 *     // You can access inject object via $.
 *     // Named options are accessible via opts.
 *     // Enumerable arguments are accessible via args.
 *     // Returned value is sent to the caller in JSON format.
 *     func1: ({$, ...opts}, ...args) => [opts, ...args],
 *
 *     // You can access connection object to operate on live data stream.
 *     func2: async ({$: {conn, state}}, initial, max) => {
 *       state.counter = initial
 *       // promise + decode will stop the handler by throwing an EOF in case the client closes the connection.
 *       new Promise(() => conn.decode()).finally()
 *       while (!conn.done && state.counter <= max) {
 *         const msg = await conn.decode()
 *         await log(msg)
 *         await conn.encode(state.counter++)
 *         await sleep(1)
 *       }
 *     },
 *   }
 * }).catch(log)
 *
 * // You can build API client by attaching target id or alias and binding api scheme to the rpc client.
 * const api = rpc.target(id).bind("foo", {"bar.baz": ["fiz"]}, "flow")
 * (async () => {
 *   // You can call api method for a single value.
 *   const item = await api.foo()
 *
 *   // You can pass arguments if needed.
 *   await api.fiz({named: "param"}, "yolo", 1, true)
 *
 *   // You can collect data asynchronously until EOF and return them as a list.
 *   const items = await api.flow
 *     .collect("pass", {some: "params"}, "if needed")
 *
 *   // You can use map operator to access or process incoming values.
 *   const mappedItems = await api.flow.map(i => {
 *     log(i) // do whatever operation
 *     return "item:" + i // optionally return mapped value
 *   }).collect()
 *
 *   // You can close process with error code.
 *   exit(1)
 * }).catch(log)
 * connect().catch(log)
 */
export const rpc = new RpcClient();
