export function bind(caller, routes) {
  const r = prepare(routes)
  const copy = caller.copy()
  for (let [method, port] of r) {
    if (caller[method]) {
      throw `method '${method}' already exist`
    }
    copy[method] = caller.call(port)
  }
  return copy
}

const prefix = /^\*/

function prepare(routes) {
  if (!Array.isArray(routes)) throw `cannot prepare routes of type ${typeof routes}`
  const prepared = []
  for (let key in routes) {
    const route = routes[key]
    switch (typeof route) {
      case "string":
        const method = route.replace(prefix, '')
        prepared.push([method, method])
        continue
      case "object":
        for (let port in route) {
          for (let method of route[port]) {
            method = method.replace(prefix, '')
            const route = [port, method].join('.')
            prepared.push([method, route])
          }
        }
    }
  }
  return prepared
}