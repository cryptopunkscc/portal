var portal = (function (exports) {
  'use strict';

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
    wails = window['go']['main']['Adapter'];
  } catch {
  }

  /* eslint-disable */
  const platform$1 = wails ? "wails" : undefined;

  /* eslint-disable */
  const adapter = () => ({
    // apphost
    astral_conn_accept: wails['ConnAccept'],
    astral_conn_close: wails['ConnClose'],
    astral_conn_read: async (id, buffer) => {
      const base64 = await wails['ConnRead'](id, buffer.byteLength);
      const binary = atob(base64);
      const len = binary.length;
      const view = new Uint8Array(buffer);
      for (let i = 0; i < len; i++) {
        view[i] = binary.charCodeAt(i);
      }
      return len
    },
    astral_conn_write: wails['ConnWrite'],
    astral_conn_read_ln: wails['ConnReadLn'],
    astral_conn_write_ln: wails['ConnWriteLn'],
    astral_node_info: wails['NodeInfo'],
    astral_query: wails['Query'],
    astral_resolve: wails['Resolve'],
    astral_service_close: wails['ServiceClose'],
    astral_service_register: wails['ServiceRegister'],
    astral_interrupt: wails['Interrupt'],
    // runtime
    sleep: wails['Sleep'],
    log: wails['Log'],
    exit: wails['Exit'],
  });

  inject(platform$1, adapter);

  // ================== Object oriented adapter ==================

  class ApphostClient {
    async register() {
      await bindings.astral_service_register();
      return new AppHostListener()
    }

    async query(target, query) {
      const json = await bindings.astral_query(target, query);
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
    constructor() {
    }

    async accept() {
      const json = await bindings.astral_conn_accept();
      const data = JSON.parse(json);
      return new ApphostConn(data)
    }

    async close() {
      await bindings.astral_service_close();
    }
  }

  class ApphostConn {
    constructor(data) {
      this.id = data.id;
      this.query = data.query;
      this.remoteId = data.remoteId;
    }

    async read(buffer) {
      try {
        return await bindings.astral_conn_read(this.id, buffer)
      } catch (e) {
        this.done = true;
        if (e === "EOF") {
          return -1
        }
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

    async readLn() {
      try {
        return await bindings.astral_conn_read_ln(this.id)
      } catch (e) {
        this.done = true;
        throw e
      }
    }

    async writeLn(data) {
      try {
        return await bindings.astral_conn_write_ln(this.id, data)
      } catch (e) {
        this.done = true;
        throw e
      }
    }

    async close() {
      if (!this.done) {
        this.done = true;
        await bindings.astral_conn_close(this.id);
      }
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

  function hasParams(query) {
    return query.search(/[?]/) > -1
  }

  /**
   * Formats an array of query parameters into a single query string.
   *
   * @param {Array} params - An array of query parameter objects and primitive type to format.
   * @return {string} A formatted query string by joining all formatted parameters with '&'.
   * @throws {TypeError} If the input is not an array.
   */
  function formatQueryParams(params) {
    if (!Array.isArray(params)) throw new TypeError('Expected an array of parameters.');
    return params.map(formatQueryParam).join('&')
  }

  function formatQueryParam(param) {
    if (param === null) return `_=null`
    if (param === undefined) return `_=undefined`
    if (!param) return `_=${encodeURIComponent(param)}`
    if (Array.isArray(param)) throw new TypeError('Expected a non-array.');
    if (typeof param === 'object') return Object.entries(param).map(e =>
      e.map(encodeURIComponent).join('=')
    ).join('&')

    return `_=${encodeURIComponent(param)}`
  }

  /**
   * Parses a query string into an object where keys map to their corresponding values.
   *
   * @param {string} query - The query string to be parsed. It should be in the format of key=value pairs separated by '&'.
   * @return {Object} - An object representing the parsed query parameters. Keys are strings, and values are strings or arrays of strings for repeated keys.
   */
  function parseQueryParams(query) {
    if (typeof query !== 'string') throw new TypeError('Expected a string.');
    let acc = {};
    query.split('&').map(parseQueryParam).forEach(([key, value]) => {
      if (key in acc) {
        acc[key] = Array.isArray(acc[key]) ? acc[key].concat(value) : [acc[key], value];
      } else if (key === '_') {
        acc[key] = [value];
      } else {
        acc[key] = value;
      }
    });
    return acc
  }

  function parseQueryParam(param) {
    if (typeof param !== 'string') throw new TypeError('Expected a string.');
    let [key, value] = param.split('=');
    key = decodeURIComponent(key);
    value = decodeURIComponent(value);
    value = parseToPrimitive(value);
    return [key, value];
  }

  function parseToPrimitive(value) {
    if (value === null || value === undefined) return value;
    if (value === "") return value;

    const num = Number(value);
    if (!isNaN(num)) return num;

    const lower = value.toLowerCase();
    if (lower === "true") return true;
    if (lower === "false") return false;
    if (lower === "null") return null;
    if (lower === "undefined") return undefined;

    return value;

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
        cmd += formatQueryParams(params);
      }
      if (cmd) await this.writeLn(cmd);
      return this
    }

    async encode(data) {
      let json = JSON.stringify(data);
      if (json === undefined) json = '{}';
      return await super.writeLn(json)
    }

    async decode() {
      const resp = await this.readLn();
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

  async function serve(client, ctx) {
    const listener = await client.register();
    listen(ctx, listener).catch(bindings.log);
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
        if (conn.done) {
          return
        }
        await conn.encode(result);
        handle = handlers;
      }
      params = await conn.readLn();
      if (typeof handle === "object") {
        [handle, params] = unfold(handle, params);
      }
    }
  }

  async function invoke(ctx, handle, params) {
    const type = typeof handle;
    switch (type) {
      case "function":
        if (!params) return await handle({$:ctx})
        const [opts, args] = preparePayload(ctx, params);
        return await handle(opts, ...args)

      case "object":
        return // skip nested router

      default:
        throw `invalid handler type ${type}`
    }
  }

  function preparePayload(ctx, params) {
    const opts = parseQueryParams(params);
    const args = opts._ ? opts._ : [];
    delete opts._;
    opts.$ = ctx;
    return [opts, args]
  }

  function unfold(handlers, query) {
    if (!query) return [handlers, query]
    let [service, args] = query.split("?");
    let chunks = service.split(".");

    for (const chunk of chunks) {
      handlers = handlers[chunk];
      if (typeof handlers === "undefined") {
        throw `cannot find handler for ${query}`
      }
    }
    return [handlers, args]
  }

  class RpcClient extends ApphostClient {

    bind(...routes) {
      return bind(this, routes);
    }

    copy(data) {
      return Object.assign(new RpcClient(), {...this, ...data});
    }

    target(id) {
      return this.copy({targetId: id})
    }

    call(port, ...params) {
      port = port ? port : "";
      return call(this, port, ...params);
    }

    async conn(port, ...params) {
      port = port ? port : "";
      const query = formatQuery(port, params);
      const conn = await super.query(this.targetId, query);
      return new RpcConn(conn)
    }

    async serve(ctx) {
      return await serve(this, ctx);
    }
  }

  function formatQuery(port, params) {
    let query = port;
    if (params.length > 0) {
      query += '?' + formatQueryParams(params);
    }
    return query
  }

  const log = async any => await bindings.log(typeof any == 'object' ? JSON.stringify(any) : any);
  const {exit, sleep, platform} = bindings;
  const apphost = new ApphostClient();
  const rpc = new RpcClient();

  exports.apphost = apphost;
  exports.exit = exit;
  exports.log = log;
  exports.platform = platform;
  exports.rpc = rpc;
  exports.sleep = sleep;

  return exports;

})({});
