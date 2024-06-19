import {bindings} from "../bindings.js";
import {RpcConn} from "./conn.js";
import {prepareRoutes} from "./route.js";
import {splitQuery} from "./query";

export async function serve(client, ctx) {
  const routes = prepareRoutes(ctx)
  for (let route of routes) {
    const listener = await client.register(route)
    listen(ctx, listener).catch(bindings.log)
  }
}

async function listen(ctx, listener) {
  for (; ;) {
    let conn = await listener.accept()
    conn = new RpcConn(conn)
    handle(ctx, conn).catch(bindings.log).finally(() =>
      conn.close().catch(bindings.log))
  }
}

async function handle(ctx, conn) {
  const inject = {...ctx.handlers, ...ctx.inject, conn: conn}
  const query = conn.query
  let [handlers, params] = unfold(ctx.handlers, query)
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
        result = await invoke(inject, handle, params)
      } catch (e) {
        result = {error: e}
      }
      await conn.encode(result)
      handle = handlers
    }
    params = await conn.read();
    if (typeof handle === "object") {
      [handle, params] = unfold(handle, params)
    }
  }
}

async function invoke(ctx, handle, params) {
  const type = typeof handle
  switch (type) {
    case "function":
      if (!params) return await handle(ctx)
      const args = JSON.parse(params)
      if (Array.isArray(args)) return await handle(...args, ctx)
      return await handle(args, ctx)

    case "object":
      return // skip nested router

    default:
      throw `invalid handler type ${type}`

  }
}

function unfold(handlers, query) {
  if (query === "") {
    return [handlers]
  }
  const [next, rest] = splitQuery(query)
  const nested = handlers[next]
  if (rest === undefined) {
    return [nested]
  }
  if (typeof nested !== "undefined") {
    return unfold(nested, rest)
  }
  if (typeof handlers === "function") {
    return [handlers, rest]
  }
  throw "cannot unfold"
}
