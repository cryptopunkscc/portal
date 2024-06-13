import {bindings} from "../bindings.js";
import {RpcConn} from "./conn.js";
import {prepareRoutes} from "./route.js";

const log = bindings.log

export async function serve(client, ctx) {
  const routes = prepareRoutes(ctx)
  for (let route of routes) {
    const listener = await client.register(route)
    listen(ctx, listener).catch(log)
  }
}

async function listen(ctx, listener) {
  for (; ;) {
    let conn = await listener.accept()
    conn = new RpcConn(conn)
    handle(ctx, conn).catch(log).finally(() =>
      conn.close().catch(log))
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
      await conn.writeJson({error: `no handler for query ${params} ${typeof handle}`})
      return
    }
    if (params || canInvoke) {
      try {
        result = await invoke(inject, handle, params)
      } catch (e) {
        result = {error: e}
      }
      await conn.writeJson(result)
      handle = handlers
    }
    params = await conn.read();
    if (typeof handle === "object") {
      [handle, params] = unfold(handle, params)
    }
  }
}

async function invoke(ctx, handle, params) {
  if (handle === undefined) {
    throw "undefined handler"
  }
  switch (typeof handle) {
    case "function":
      const args = JSON.parse(params)
      if (Array.isArray(args)) {
        return await handle(...args, ctx)
      }
      return await handle(args, ctx)
    case "object":
      return
  }
}

function unfold(handlers, query) {
  if (query === "") {
    return [handlers]
  }
  const [next, rest] = split(query)
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

function split(query) {
  const index = query.search(/[?.{\[]/)
  if (index === -1) {
    return [query]
  }
  const left = query.slice(0, index)
  let right = query.slice(index, query.length)
  if (/^[.?]/.test(right)) {
    right = right.slice(1)
  }
  return [left, right]
}
