import {bindings} from "../bindings.js";
import {RpcConn} from "./conn.js";
import {parseQueryParams} from "./params";

/**
 * Registers given API on {@link RpcClient}.
 *
 * @async
 * @param {RpcClient} client
 * @param {object} api
 * @returns {Promise<void>}
 */
export async function serve(client, api) {
  const listener = await client.register()
  listen(api, listener).catch(bindings.log)
}

/**
 * Accepts incoming connections and handles them in the context of given API.
 *
 * @async
 * @param {object} api
 * @param {ApphostListener} listener
 * @returns {Promise<void>}
 */
async function listen(api, listener) {
  for (; ;) {
    let conn = await listener.accept()
    conn = new RpcConn(conn)
    handle(api, conn).catch(bindings.log).finally(() =>
      conn.close().catch(bindings.log))
  }
}

/**
 * Handles incoming {@link RpcConn} in the context of given API.
 *
 * @async
 * @param {object} api
 * @param {RpcConn} conn
 * @returns {Promise<void>}
 */
async function handle(api, conn) {
  const inject = {...api.handlers, ...api.inject, conn: conn}
  const query = conn.query
  let [handlers, params] = unfold(api.handlers, query)
  let handle = handlers
  let result
  let canInvoke
  for (; ;) {
    canInvoke = typeof handle === "function"
    if (params && !canInvoke) {
      await conn.encode({error: `no handler for query ${params} ${typeof handle}`})
      return
    }
    if (params || canInvoke) {
      try {
        result = await invoke(handle, inject, params)
      } catch (e) {
        result = {error: e}
      }
      if (conn.done) {
        return
      }
      await conn.encode(result)
      handle = handlers
    }
    params = await conn.readLn();
    if (typeof handle === "object") {
      [handle, params] = unfold(handle, params)
    }
  }
}


/**
 * Invokes handler with given api and params.
 *
 * @async
 * @param {function} handle
 * @param {(object|string)[]} api
 * @param {any[]} params
 * @returns {Promise<*>}
 */
async function invoke(handle, api, params) {
  const type = typeof handle
  switch (type) {
    case "function":
      if (!params) return await handle({$: api})
      const [opts, args] = preparePayload(api, params)
      return await handle(opts, ...args)

    case "object":
      return // skip nested router

    default:
      throw `invalid handler type ${type}`
  }
}

/**
 * Prepares options and arguments for handle.
 *
 * @param api
 * @param params
 * @returns {[any[], any[]]}
 */
function preparePayload(api, params) {
  const opts = parseQueryParams(params)
  const args = opts._ ? opts._ : []
  delete opts._
  opts.$ = api
  return [opts, args]
}

/**
 * Parses query and returns its arguments attached to the corresponding handler(s).
 *
 * @param {object} handlers
 * @param {string} query
 * @returns {[function|object, any[]]}
 */
function unfold(handlers, query) {
  if (!query) return [handlers, query]
  let [service, args] = query.split("?")
  let chunks = service.split(".")

  for (const chunk of chunks) {
    handlers = handlers[chunk]
    if (typeof handlers === "undefined") {
      throw `cannot find handler for ${query}`
    }
  }
  return [handlers, args]
}
