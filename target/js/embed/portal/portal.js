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


let wails = {};
try {
  wails = window;
} catch (e) {

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
  log: async (...arg) => await wails['go']['main']['Adapter']['LogArr'](arg),
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
    log: (arg1) => _app_host.logArr(JSON.stringify(arg1)),
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

class RpcConn extends ApphostConn {
  constructor(data) {
    super(data);
  }

  async encode(data) {
    let json = JSON.stringify(data);
    return await super.write(json + '\n')
  }

  async decode() {
    const resp = await this.read();
    return JSON.parse(resp)
  }

  async call(method, ...params) {
    let cmd = method;
    if (params) {
      cmd += '?' + JSON.stringify(params);
    }
    await this.write(cmd + '\n');
  }

  async request(query, ...params) {
    await this.call(query, ...params);
    return await this.decode()
  }

  async observe(consume) {
    for (;;) {
      const next = await this.decode();
      const last = await consume(next);
      this.value = next;
      if (last) {
        return last
      }
    }
  }

  caller(method) {
    return async (...params) => await this.call(method, ...params)
  }

  requester(method) {
    return async (...params) => await this.request(method, ...params)
  }

  observer(method) {
    return (...params) => ({
      observe: async (consume) => {
        await this.call(method, ...params);
        return await this.observe(consume)
          .finally(() => this.close()) // TODO consider if it's ok to close conn automatically.
      }
    })
  }

  bind(...methods) {
    for (let method of methods) {
      const collect = method[0] === '*';
      if (collect) {
        method = method.split(1);
      }
      if (this[method]) {
        throw `method '${method}' already exist`
      }
      if (collect) {
        this[method] = this.observer(method);
      } else {
        this[method] = this.requester(method);
      }
    }
    return this
  }
}

function prepareRoutes$1(ctx) {
  let routes = resolveRoutes(ctx.handlers);
  routes = formatRoutes$1(routes);
  routes = maskRoutes$1(routes, ctx.routes);
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

function formatRoutes$1(routes) {
  const formatted = [];
  for (let route of routes) {
    formatted.push(route.join("."));
  }
  return formatted
}

function maskRoutes$1(routes, masks) {
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
  const routes = prepareRoutes$1(ctx);
  for (let route of routes) {
    const listener = await client.register(route);
    listen$1(ctx, listener).catch(bindings.log);
  }
}

async function listen$1(ctx, listener) {
  for (; ;) {
    let conn = await listener.accept();
    conn = new RpcConn(conn);
    handle$1(ctx, conn).catch(bindings.log).finally(() =>
      conn.close().catch(bindings.log));
  }
}

async function handle$1(ctx, conn) {
  const inject = {...ctx.handlers, ...ctx.inject, conn: conn};
  const query = conn.query;
  let [handlers, params] = unfold$1(ctx.handlers, query);
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
        result = await invoke$1(inject, handle, params);
      } catch (e) {
        result = {error: e};
      }
      await conn.encode(result);
      handle = handlers;
    }
    params = await conn.read();
    if (typeof handle === "object") {
      [handle, params] = unfold$1(handle, params);
    }
  }
}

async function invoke$1(ctx, handle, params) {
  if (!handle) {
    throw "undefined handler"
  }
  switch (typeof handle) {
    case "function":
      if (!params) {
        return await handle(ctx)
      }

      const args = JSON.parse(params);
      if (Array.isArray(args)) {
        return await handle(...args, ctx)
      }
      return await handle(args, ctx)
    case "object":
      return
  }
}

function unfold$1(handlers, query) {
  if (query === "") {
    return [handlers]
  }
  const [next, rest] = split$1(query);
  const nested = handlers[next];
  if (rest === undefined) {
    return [nested]
  }
  if (typeof nested !== "undefined") {
    return unfold$1(nested, rest)
  }
  if (typeof handlers === "function") {
    return [handlers, rest]
  }
  throw "cannot unfold"
}

function split$1(query) {
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

class RpcClient extends ApphostClient {

  constructor() {
    super();
    this.port = "";
  }

  // TODO deprecated use call instead
  async query(query) {
    let conn = await super.query(query, this.targetId);
    conn = new RpcConn(conn);
    return conn
  }

  async serve(ctx) {
    await serve(this, ctx);
  }

  async call(query, ...params) {
    if (params.length > 0) {
      query += '?' + JSON.stringify(params);
    }
    const conn = await super.query(query, this.targetId);
    return new RpcConn(conn)
  }

  async request(query, ...params) {
    const conn = await this.call(query, ...params);
    const response = await conn.decode();
    conn.close().catch(bindings.log);
    return response
  }

  caller(query) {
    return async (...params) => await this.call(query, ...params)
  }

  requester(query) {
    return async (...params) => await this.request(query, ...params)
  }

  observer(query) {
    return async (...params) => {
      if (typeof params[params.length - 1] !== "function") {
        return await this.request(query, ...params)
      }
      const consume = params.pop();
      const conn = await this.call(query, ...params);
      let last;
      try {
        last = await conn.observe(consume);
      } catch (e) {
        if ('{}' === JSON.stringify(e)) {
          last = conn.value;
        } else {
          throw e
        }
      } finally {
        conn.close().finally();
      }
      return last
    }
  }

  copy(data) {
    return Object.assign(new RpcClient(), {...this, ...data})
  }

  target(id) {
    return this.copy({targetId: id})
  }

  bind(route, ...methods) {
    const copy = this.copy();
    for (let method of methods) {
      const collect = method[0] === '*';
      if (collect) {
        method = method.substring(1); // drop * prefix
      }
      if (this[method]) {
        throw `method '${method}' already exist`
      }
      const port = [route, method].join('.');
      if (collect) {
        copy[method] = copy.observer(port);
      } else {
        copy[method] = copy.requester(port);
      }
    }
    return copy
  }
}

ApphostConn.prototype.readJson = async function (method) {
  const resp = await this.read();
  const json = JSON.parse(resp);
  return json
};

ApphostConn.prototype.jsonReader = async function (method) {
  const read = async () => await this.readJson(method);
  read.cancel = async () => await this.close();
  return read
};

ApphostConn.prototype.writeJson = async function (data) {
  // if (Array.isArray(data) && data.length === 1) {
  //   data = data[0]
  // }
  const json = JSON.stringify(data);
  // log(this.id + " conn => " + this.query + ":" + json.trimEnd())
  await this.write(json + '\n');
};

ApphostConn.prototype.rpcCall = async function (method, ...data) {
  let cmd = method;
  if (data.length > 0) {
    cmd += "?" + JSON.stringify(data);
  }
  // log(this.id + " conn => " + this.query + "." + cmd)
  await this.write(cmd + '\n');
};

ApphostConn.prototype.rpcQuery = function (method) {
  const conn = this;
  return async function (...data) {
    // log("conn rpc query", method)
    await conn.rpcCall(method, ...data);
    return await conn.readJson(method)
  }
};

// Bind RPC api of service associated to this connection
ApphostConn.prototype.bindRpc = async function () {
  const conn = this;
  // request api methods
  await conn.rpcCall("api");

  // read api methods
  const methods = await conn.readJson("api");

  // bind methods
  for (let method of methods) {
    conn[method] = async (...data) => {
      await conn.rpcCall(method, ...data);
      return await conn.readJson(method)
    };
  }

  // bind subscribe
  conn.subscribe = async (method, ...data) => {
    await conn.rpcCall(method, ...data);
    return conn.jsonReader(method)
  };
};

const {log: log$3} = bindings;

ApphostClient.prototype.rpcCall = async function (identity, service, method, ...data) {
  let cmd = service;
  if (method) {
    cmd += "." + method;
  }
  if (data.length > 0) {
    cmd += "?" + JSON.stringify(data);
  }
  const conn = await this.query(cmd, identity);
  // log(conn.id + " client => " + cmd)
  return conn
};

ApphostClient.prototype.rpcQuery = function (port, identity) {
  const client = this;
  return async function (...data) {
    const conn = await client.rpcCall(identity, port, "", ...data);
    const json = await conn.readJson(port);
    conn.close().catch(log$3);
    return json
  }
};

ApphostClient.prototype.bindRpc = async function (identity, service) {
  const client = this;
  // request api methods
  const conn = await client.rpcCall(identity, service, "api");

  // read api methods
  const methods = await conn.readJson("api");
  conn.close().catch(log$3);

  // bind methods
  for (let method of methods) {
    client[method] = async (...data) => {
      const conn = await client.rpcCall(identity, service, method, ...data);
      const json = await conn.readJson(method);
      conn.close().catch(log$3);
      return json
    };
  }

  // bind subscribe
  client.subscribe = async (method, ...data) => {
    const conn = await client.rpcCall(identity, service, method, ...data);
    return await conn.jsonReader(method)
  };
  return client
};

const {log: log$2} = bindings;


// Bind RPC service to given name
ApphostClient.prototype.bindRpcService = async function (Service) {
  const props = Object.getOwnPropertyNames(Service.prototype);
  if (props[0] !== "constructor") throw new Error("Service must have a constructor")
  const methods = props.slice(1, props.length);
  methods.push("api");
  Service.prototype.api = async () => {
    return methods
  };
  const srv = new Service();
  const listener = await this.register(srv.name + "*");
  // log("listen " + srv.name)
  astral_rpc_listen(srv, listener).catch(log$2);
  return listener
};

async function astral_rpc_listen(srv, listener) {
  for (; ;) {
    const conn = await listener.accept();
    // log(conn.id + " service <= " + conn.query)
    try {
      astral_rpc_handle(srv, conn).catch(log$2);
    } catch (e) {
      // log(conn.id + " service !! " + conn.query + ":" + e)
      conn.close().catch(log$2);
    }
  }
}

async function astral_rpc_handle(srv, conn) {
  let query = conn.query.slice(srv.name.length);
  let method = query, args = [];
  const single = query !== '';
  const write = async (data) => await conn.writeJson(data);
  const read = async (method) => await conn.readJson(method);

  for (; ;) {
    if (!single) {
      query = await conn.read();
      // log(conn.id + " service <== " + query)
    }
    [method, args] = parseQuery(query);

    let result;
    try {
      result = await srv[method](...args, write, read);
    } catch (e) {
      result = {error: e};
    }
    if (result !== undefined) {
      await conn.writeJson(result);
    }
    if (single) {
      conn.close().catch(log$2);
      break
    }
  }
}

function parseQuery(query) {
  if (query[0] === '.') {
    query = query.slice(1);
  }
  const match = /[?\[{]/.exec(query);
  const method = query.slice(0, match.index);
  let payload = query.slice(match.index);
  if (payload[0] === '?') {
    payload = payload.slice(1);
  }
  let args = [];
  if (payload) {
    args = JSON.parse(payload);
  }
  return [method, args]
}

const log$1 = bindings.log;


ApphostClient.prototype.registerRpc = async function (ctx) {
  const routes = prepareRoutes(ctx);
  for (let route of routes) {
    const listener = await this.register(route);
    listen(ctx, listener).catch(log$1);
  }
};

function prepareRoutes(ctx) {
  let routes = collectRoutes(ctx.handlers);
  routes = formatRoutes(routes);
  routes = maskRoutes(routes, ctx.routes);
  return routes
}

function collectRoutes(handlers, ...name) {
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
    const nested = collectRoutes(next, ...[...name, prop]);
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

async function listen(ctx, listener) {
  for (; ;) {
    const conn = await listener.accept();
    try {
      handle(ctx, conn).catch(log$1);
    } catch (e) {
      conn.close().catch(log$1);
    }
  }
}

async function handle(ctx, conn) {
  const inject = {...ctx.handlers, ...ctx.inject, conn: conn};
  let [handlers, params] = unfold(ctx.handlers, conn.query);
  let handle = handlers;
  let result;
  let canInvoke;
  for (; ;) {
    canInvoke = typeof handle === "function";
    if (params && !canInvoke) {
      await conn.writeJson({error: `no handler for query ${params} ${typeof handle}`});
      return
    }
    if (params || canInvoke) {
      try {
        result = await invoke(inject, handle, params);
      } catch (e) {
        result = {error: e};
      }
      await conn.writeJson(result);
      handle = handlers;
    }
    params = await conn.read();
    if (typeof handle === "object") {
      [handle, params] = unfold(handle, params);
    }
  }
}

async function invoke(ctx, handle, params) {
  if (handle === undefined) {
    throw "undefined handler"
  }
  switch (typeof handle) {
    case "function":
      const args = JSON.parse(params);
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
  const [next, rest] = split(query);
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

function split(query) {
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

const {log, sleep, platform} = bindings;
const apphost = new ApphostClient();
const rpc = new RpcClient();

export { apphost, log, platform, rpc, sleep };
