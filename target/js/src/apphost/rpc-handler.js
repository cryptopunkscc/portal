import {ApphostClient} from "./adapter.js";
import {bindings} from "../bindings.js";
const log = bindings.log


ApphostClient.prototype.registerRpc = async function (ctx) {
  const routes = prepareRoutes(ctx)
  for (let route of routes) {
    const listener = await this.register(route)
    listen(ctx, listener).catch(log)
  }
}

function prepareRoutes(ctx) {
  let routes = collectRoutes(ctx.handlers)
  routes = formatRoutes(routes)
  routes = maskRoutes(routes, ctx.routes)
  return routes
}

function collectRoutes(handlers, ...name) {
  if (typeof handlers !== "object") {
    return name
  }

  const props = Object.getOwnPropertyNames(handlers)
  if (props.length === 0) {
    return name
  }
  const routes = []
  for (let prop of props) {
    const next = handlers[prop]
    const nested = collectRoutes(next, ...[...name, prop])
    if (typeof nested[0] === "string") {
      routes.push(nested)
    } else {
      routes.push(...nested)
    }
  }
  return routes
}

function formatRoutes(routes) {
  const formatted = []
  for (let route of routes) {
    formatted.push(route.join("."))
  }
  return formatted
}

function maskRoutes(routes, masks) {
  masks = masks ? masks : []
  let arr = [...routes]
  for (let mask of masks) {
    const last = mask.length - 1;
    if (/[*:]/.test(mask.slice(last))) {
      mask = mask.slice(0, last)
    }
    arr = arr.filter(val => !val.startsWith(mask))
  }
  masks = masks.filter(mask => !mask.endsWith(":"))
  arr.push(...masks)
  return arr
}

async function listen(ctx, listener) {
  for (; ;) {
    const conn = await listener.accept()
    try {
      handle(ctx, conn).catch(log)
    } catch (e) {
      conn.close().catch(log)
    }
  }
}

async function handle(ctx, conn) {
  const inject = {...ctx.handlers, ...ctx.inject, conn: conn}
  let [handlers, params] = unfold(ctx.handlers, conn.query)
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
      } else {
        return await handle(args, ctx)
      }
    case "object":
      return
  }
}

function unfold(handlers, query) {
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
