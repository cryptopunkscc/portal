export function prepareRoutes(ctx) {
  let routes = resolveRoutes(ctx.handlers)
  routes = formatRoutes(routes)
  routes = maskRoutes(routes, ctx.routes)
  return routes
}

function resolveRoutes(handlers, ...name) {
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
    const nested = resolveRoutes(next, ...[...name, prop])
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
    if (mask === '*') {
      return [masks]
    }
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
