const bindings = {};

function inject(platform, adapter) {
  if (platform !== undefined) {
    Object.assign(bindings, {
      platform: platform,
      ...adapter()
    });
  }
}

// ================== Wails bindings adapter ==================


let wails;
try {
  wails = window;
} catch (e) {
  wails = {};
}

/* eslint-disable */
const platform$3 = typeof wails['go'] === "undefined" ? undefined : "wails";

/* eslint-disable */
const adapter$2 = () => ({
  // apphost
  astral_conn_accept: wails['go']['main']['Adapter']['ConnAccept'],
  astral_conn_close: wails['go']['main']['Adapter']['ConnClose'],
  astral_conn_read: wails['go']['main']['Adapter']['ConnRead'],
  astral_conn_write: wails['go']['main']['Adapter']['ConnWrite'],
  astral_node_info: wails['go']['main']['Adapter']['NodeInfo'],
  astral_query: wails['go']['main']['Adapter']['Query'],
  astral_query_name: wails['go']['main']['Adapter']['QueryName'],
  astral_resolve: wails['go']['main']['Adapter']['Resolve'],
  astral_service_close: wails['go']['main']['Adapter']['ServiceClose'],
  astral_service_register: wails['go']['main']['Adapter']['ServiceRegister'],
  astral_interrupt: wails['go']['main']['Adapter']['Interrupt'],
  // runtime
  sleep: wails['go']['main']['Adapter']['Sleep'],
  log: wails['go']['main']['Adapter']['Log'],
});

inject(platform$3, adapter$2);

// ================== Android bindings adapter ==================

/* eslint-disable */
const platform$2 = typeof _app_host === "undefined" ? undefined : "android";

/* eslint-disable */
const adapter$1 = () => {

  const _awaiting = new Map();

  window._resolve = (id, value) => {
    _awaiting.get(id)[0](value);
    _awaiting.delete(id);
  };

  window._reject = (id, error) => {
    _awaiting.get(id)[1](error);
    _awaiting.delete(id);
  };

  const _promise = (block) =>
    new Promise((resolve, reject) =>
      _awaiting.set(block(), [resolve, reject]));

  return {
    // apphost
    astral_node_info: (arg1) => _promise(() => _app_host.nodeInfo(arg1)).then(v => JSON.parse(v)),
    astral_conn_accept: (arg1) => _promise(() => _app_host.connAccept(arg1)),
    astral_conn_close: (arg1) => _promise(() => _app_host.connClose(arg1)),
    astral_conn_read: (arg1) => _promise(() => _app_host.connRead(arg1)),
    astral_conn_write: (arg1, arg2) => _promise(() => _app_host.connWrite(arg1, arg2)),
    astral_query: (arg1, arg2) => _promise(() => _app_host.query(arg1, arg2)),
    astral_query_name: (arg1, arg2) => _promise(() => _app_host.queryName(arg1, arg2)),
    astral_resolve: (arg1) => _promise(() => _app_host.resolve(arg1)),
    astral_service_close: (arg1) => _promise(() => _app_host.serviceClose(arg1)),
    astral_service_register: (arg1) => _promise(() => _app_host.serviceRegister(arg1)),
    astral_interrupt: () => _promise(() => _app_host.interrupt()),
    // runtime
    sleep: (arg1) => _promise(() => _app_host.sleep(arg1)),
    log: (arg1) => _app_host.log(arg1),
  }
};

inject(platform$2, adapter$1);

const platform$1 = typeof _log === 'undefined' ? undefined : "common";

/* eslint-disable */
const adapter = () => ({
  // apphost
  astral_conn_accept: _astral_conn_accept,
  astral_conn_close: _astral_conn_close,
  astral_conn_read: _astral_conn_read,
  astral_conn_write: _astral_conn_write,
  astral_node_info: _astral_node_info,
  astral_query: _astral_query,
  astral_query_name: _astral_query_name,
  astral_resolve: _astral_resolve,
  astral_service_close: _astral_service_close,
  astral_service_register: _astral_service_register,
  astral_interrupt: _astral_interrupt,
  // apphost
  sleep: _sleep,
  log: _log,
});

inject(platform$1, adapter);

// ================== Object oriented adapter ==================

class ApphostClient {
  async register(service) {
    await bindings.astral_service_register(service);
    return new AppHostListener(service)
  }

  async query(query, identity) {
    identity = identity ? identity : "";
    const json = await bindings.astral_query(identity, query);
    const data = JSON.parse(json);
    return new ApphostConn(data, query)
  }

  async queryName(name, query) {
    const json = await bindings.astral_query_name(name, query);
    const data = JSON.parse(json);
    return new ApphostConn(data, query)
  }

  async nodeInfo(id) {
    return await bindings.astral_node_info(id)
  }

  async resolve(name) {
    return await bindings.astral_resolve(name)
  }

  async interrupt() {
    await bindings.astral_interrupt();
  }
}

class AppHostListener {
  constructor(port) {
    this.port = port;
  }

  async accept() {
    const json = await bindings.astral_conn_accept(this.port);
    const data = JSON.parse(json);
    return new ApphostConn(data)
  }

  async close() {
    await bindings.astral_service_close(this.port);
  }
}

class ApphostConn {
  constructor(data) {
    this.id = data.id;
    this.query = data.query;
    this.remoteId = data.remoteId;
  }

  async read() {
    try {
      return await bindings.astral_conn_read(this.id)
    } catch (e) {
      this.done = true;
      throw e
    }
  }

  async write(data) {
    try {
      return await bindings.astral_conn_write(this.id, data)
    } catch (e) {
      this.done = true;
      throw e
    }
  }

  async close() {
    this.done = true;
    await bindings.astral_conn_close(this.id);
  }
}

function bind(caller, routes) {
  const r = prepare(routes);
  const copy = caller.copy();
  for (let [method, port] of r) {
    if (copy[method]) {
      throw `method '${method}' already exist`
    }
    copy[method] = copy.call(port);
  }
  return copy
}

const prefix = /^\*/;

function prepare(routes) {
  if (!Array.isArray(routes)) throw `cannot prepare routes of type ${typeof routes}`
  const prepared = [];
  for (let key in routes) {
    const route = routes[key];
    switch (typeof route) {
      case "string":
        const method = route.replace(prefix, '');
        prepared.push([method, method]);
        continue
      case "object":
        for (let port in route) {
          for (let method of route[port]) {
            method = method.replace(prefix, '');
            const route = [port, method].join('.');
            prepared.push([method, route]);
          }
        }
    }
  }
  return prepared
}

/**
 * @param {RpcClient | RpcConn} client
 * @param {string} port
 * @param {any[]} params
 */
function call(client, port, ...params) {
  const call = new RpcCall(client, port, params);
  let f = async (...params) => {
    return await call.request(...params);
  };
  return Object.assign(f, {
    inner: call,
    map: (...args) => call.map(...args),
    filter: (...args) => call.filter(...args),
    request: async (...args) => await call.request(...args),
    collect: async (...args) => await call.collect(...args),
    conn: async (...args) => await call.conn(...args),
  })
}

class RpcCall {

  mapper = arg => arg
  params = []
  single = true

  /**
   * @param {RpcClient | RpcConn} client
   * @param {string} port
   * @param {any[]} params
   */
  constructor(client, port, params) {
    this.client = client;
    this.port = port;
    this.params = Array.isArray(params) ? params : params ? [params] : [];
  }

  map(f) {
    const map = this.mapper;
    this.mapper = arg => f(map(arg));
    return this
  }

  filter(f) {
    return this.map(arg => {
      if (f(arg)) return arg
    })
  }

  async request(...params) {
    if (params.length > 0) this.params = params;
    return await this.#consume(async conn => await conn.request(...params));
  }

  async collect(...params) {
    if (params.length > 0) this.params = params;
    return await this.#consume(async conn => await conn.collect(...params));
  }

  async #consume(f) {
    const conn = await this.conn();
    conn.mapper = this.mapper;
    this.result = await f(conn);
    this.mapper = a => a; // reset mapper between requests
    if (this.single) await conn.close().catch(bindings.log);
    return this.result
  }

  async conn(...params) {
    const args = params.length > 0 ? params : this.params;
    return this.client.conn(this.port, ...args);
  }
}

function splitQuery(query) {
  const index = query.search(/[?.{\[]/);
  if (index === -1) {
    return [query]
  }
  const left = query.slice(0, index);
  let right = query.slice(index, query.length);
  if (/^[.?]/.test(right)) {
    right = right.slice(1);
  }
  return [left, right]
}


function hasParams(query) {
  return query.search(/[?{\[]/) > -1
}

class RpcConn extends ApphostConn {

  constructor(data) {
    super(data);
  }

  #sub(port) {
    if (hasParams(this.query)) throw `cannot nest connection for complete query ${chunks}`
    return port
  }

  bind(...routes) {
    return bind(this, routes);
  }

  copy() {
    return this;
  }

  call(port, ...params) {
    const c = call(this, this.#sub(port), ...params);
    c.inner.single = false;
    return c
  }

  map(f) {
    if (this.mapper) {
      const map = this.mapper;
      this.mapper = arg => f(map(arg));
    } else {
      this.mapper = f;
    }
    return this
  }

  async conn(method, ...params) {
    let cmd = method ? method : "";
    if (params.length > 0) {
      if (cmd) cmd += '?';
      cmd += JSON.stringify(params);
    }
    if (cmd) await this.write(cmd + '\n');
    return this
  }

  async encode(data) {
    let json = JSON.stringify(data);
    if (json === undefined) json = '{}';
    return await super.write(json + '\n')
  }

  async decode() {
    const resp = await this.read();
    const parsed = JSON.parse(resp);
    if (parsed === null) return null
    if (parsed.error) throw parsed.error
    return parsed
  }

  async request(...params) {
    const map = this.mapper;
    this.result = null;
    for (; ;) {
      const next = await this.decode();
      if (next === undefined) continue
      if (next === null) return this.result
      this.result = next;
      if (!map) return next
      const last = await map(next);
      if (last === undefined) continue
      if (last === null) return this.result
      return last
    }
  }

  /**
   * Collects all decoded values mapped as not null until decodes null or maps into undefined.
   *
   * @param {...any} params
   * @returns {Promise<any[]>}
   */
  async collect(...params) {
    const map = this.mapper ? this.mapper : null;
    let push;
    if (!map) push = next => this.result.push(next);
    else push = async (next) => {
      next = await map.call(this, next);
      if (next === null) return this.result
      if (next) this.result.push(next);
    };
    this.result = [];
    for (; ;) {
      let next = await this.decode();
      if (next === null) return this.result
      push(next);
    }
  }
}

function prepareRoutes(ctx) {
  let routes = resolveRoutes(ctx.handlers);
  routes = formatRoutes(routes);
  routes = maskRoutes(routes, ctx.routes);
  return routes
}

function resolveRoutes(handlers, ...name) {
  if (typeof handlers !== "object") {
    return name
  }
  const props = Object.getOwnPropertyNames(handlers);
  if (props.length === 0) {
    return name
  }
  const routes = [];
  for (let prop of props) {
    const next = handlers[prop];
    const nested = resolveRoutes(next, ...[...name, prop]);
    if (typeof nested[0] === "string") {
      routes.push(nested);
    } else {
      routes.push(...nested);
    }
  }
  return routes
}

function formatRoutes(routes) {
  const formatted = [];
  for (let route of routes) {
    formatted.push(route.join("."));
  }
  return formatted
}

function maskRoutes(routes, masks) {
  masks = masks ? masks : [];
  let arr = [...routes];
  for (let mask of masks) {
    if (mask === '*') {
      return [masks]
    }
    const last = mask.length - 1;
    if (/[*:]/.test(mask.slice(last))) {
      mask = mask.slice(0, last);
    }
    arr = arr.filter(val => !val.startsWith(mask));
  }
  masks = masks.filter(mask => !mask.endsWith(":"));
  arr.push(...masks);
  return arr
}

async function serve(client, ctx) {
  const routes = prepareRoutes(ctx);
  for (let route of routes) {
    const listener = await client.register(route);
    listen(ctx, listener).catch(bindings.log);
  }
}

async function listen(ctx, listener) {
  for (; ;) {
    let conn = await listener.accept();
    conn = new RpcConn(conn);
    handle(ctx, conn).catch(bindings.log).finally(() =>
      conn.close().catch(bindings.log));
  }
}

async function handle(ctx, conn) {
  const inject = {...ctx.handlers, ...ctx.inject, conn: conn};
  const query = conn.query;
  let [handlers, params] = unfold(ctx.handlers, query);
  let handle = handlers;
  let result;
  let canInvoke;
  for (; ;) {
    canInvoke = typeof handle === "function";
    if (params && !canInvoke) {
      await conn.encode({error: `no handler for query ${params} ${typeof handle}`});
      return
    }
    if (params || canInvoke) {
      try {
        result = await invoke(inject, handle, params);
      } catch (e) {
        result = {error: e};
      }
      await conn.encode(result);
      handle = handlers;
    }
    params = await conn.read();
    if (typeof handle === "object") {
      [handle, params] = unfold(handle, params);
    }
  }
}

async function invoke(ctx, handle, params) {
  const type = typeof handle;
  switch (type) {
    case "function":
      if (!params) return await handle(ctx)
      const args = JSON.parse(params);
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
  const [next, rest] = splitQuery(query);
  const nested = handlers[next];
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

class RpcClient extends ApphostClient {

  bind(...routes) {
    return bind(this, routes);
  }

  copy(data) {
    return Object.assign(new RpcClient(), {...this, ...data});
  }

  target(id) {
    this.targetId = id;
    return this
  }

  call(port, ...params) {
    return call(this, port, ...params);
  }

  async conn(port, ...params) {
    const query = formatQuery(port, params);
    const conn = await super.query(query, this.targetId);
    return new RpcConn(conn)
  }

  async serve(ctx) {
    return await serve(this, ctx);
  }
}

function formatQuery(port, params) {
  let query = port;
  if (params.length > 0) {
    query += '?' + JSON.stringify(params);
  }
  return query
}

const log = any => bindings.log(typeof any == 'object' ? JSON.stringify(any) : any);
const {sleep, platform} = bindings;
const apphost = new ApphostClient();
const rpc = new RpcClient();

export { apphost, log, platform, rpc, sleep };
