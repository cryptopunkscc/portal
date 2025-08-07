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

  // ================== Android bindings adapter ==================

  /* eslint-disable */
  const platform$1 = typeof _app_host === "undefined" ? undefined : "android";

  /* eslint-disable */
  const adapter = () => {

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
      astral_conn_read: (arg1, arg2) => _promise(() => {
        // TODO write result to byte array
        return _app_host.connRead(arg1, arg2);
      }),
      astral_conn_write: (arg1, arg2) => _promise(() => _app_host.connWrite(arg1, arg2)),
      astral_conn_read_ln: (arg1) => _promise(() => _app_host.connReadLn(arg1)),
      astral_conn_write_ln: (arg1, arg2) => _promise(() => _app_host.connWriteLn(arg1, arg2)),
      astral_query: (arg1, arg2) => _promise(() => _app_host.query(arg1, arg2)).then(v => JSON.parse(v)),
      astral_resolve: (arg1) => _promise(() => _app_host.resolve(arg1)),
      astral_service_close: (arg1) => _promise(() => _app_host.serviceClose(arg1)),
      astral_service_register: (arg1) => _promise(() => _app_host.serviceRegister(arg1)),
      astral_interrupt: () => _promise(() => _app_host.interrupt()),
      // runtime
      sleep: (arg1) => _promise(() => _app_host.sleep(arg1)),
      log: (arg1) => _app_host.log(arg1),
      exit: (arg1) => _app_host.exit(arg1),
    }
  };

  inject(platform$1, adapter);

  // ================== Object oriented adapter ==================

  class ApphostClient {
    async register() {
      await bindings.astral_service_register();
      return new ApphostListener()
    }

    async query(target, query) {
      const data = await bindings.astral_query(target, query);
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

  class ApphostListener {
    constructor() {
    }

    async accept() {
      const data = await bindings.astral_conn_accept();
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

  /**
   * Returns a copy of caller with each route bound as a method.
   *
   * @param {RpcClient|RpcConn} caller
   * @param {Array<string|object>} routes
   * @returns {(function&RpcClient)|(function&RpcConn)}
   */
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
   * Combines RPC client interface with its request function.
   *
   * @param {RpcClient | RpcConn} client
   * @param {string} port
   * @param {...any} params
   * @returns {function & RpcCall}
   */
  function call(client, port, ...params) {
    const call = new RpcCall(client, port, params);
    let f = call.request.bind(call);
    return Object.assign(f, {
      inner: call,
      map: (f) => call.map(f),
      filter: (f) => call.filter(f),
      request: async (...args) => await call.request(...args),
      collect: async (...args) => await call.collect(...args),
      conn: async (...args) => await call.conn(...args),
    })
  }

  /**
   * Represents a further intention for opening connections on a specific port with a given params.
   */
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

    /**
     * Adds response mapper.
     *
     * @param {(any) => any} f - mapping function
     * @returns RpcCall - new instance
     */
    map(f) {
      const map = this.mapper;
      this.mapper = arg => f(map(arg));
      return this
    }

    /**
     * Adds response filter.
     *
     * @param {(any) => boolean} f - filtering function
     * @returns RpcCall - new instance
     */
    filter(f) {
      return this.map(arg => {
        if (f(arg)) return arg
      })
    }

    /**
     * Creates new {@link RpcConn} and calls {@link RpcConn.request}.
     *
     * @async
     * @param {...any} params - connection params
     * @returns {Promise<any>}
     */
    async request(...params) {
      if (params.length > 0) this.params = params;
      return await this.#consume(async conn => await conn.request());
    }

    /**
     * Creates new {@link RpcConn} and calls {@link RpcConn.collect}.
     *
     * @async
     * @param {...any} params - connection params
     * @returns {Promise<any[]>}
     */
    async collect(...params) {
      if (params.length > 0) this.params = params;
      return await this.#consume(async conn => await conn.collect());
    }

    async #consume(f) {
      const conn = await this.conn();
      conn.mapper = this.mapper;
      this.result = await f(conn);
      this.mapper = a => a; // reset mapper between requests
      if (this.single) await conn.close().catch(bindings.log);
      return this.result
    }

    /**
     * Returns new {@link RpcConn}.
     *
     * @async
     * @param {...any} params - connection params
     * @returns {Promise<RpcConn>}
     */
    async conn(...params) {
      const args = params.length > 0 ? params : this.params;
      return this.client.conn(this.port, ...args);
    }
  }

  /**
   * Formats an array of query parameters into a single query string.
   *
   * @param {any[]} params - An array of objects and values to format as query parameters.
   * @return {string} A query parameters string generated by joining all formatted objects and values using '&'.
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
   * Parses a query string into an object, mapping keys to their corresponding values.
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

  function formatQuery(port, params) {
    let query = port;
    if (params.length > 0) {
      query += '?' + formatQueryParams(params);
    }
    return query
  }

  function hasParams(query) {
    return query.search(/[?]/) > -1
  }

  /**
   * Adds RPC implementation to {@link ApphostConn}. Compatible with astral query.
   *
   * @extends {ApphostConn}
   */
  class RpcConn extends ApphostConn {

    constructor(data) {
      super(data);
    }

    /**
     * Returns a copy of the object with each route bound as a method.
     *
     * @param {...string|object} routes
     * @return {RpcConn&object}
     *
     * @example
     * const conn = rpc.target(id).conn()
     * const api = conn.bind("foo", "bar")
     * await api.foo()
     * await api.bar("baz", 1, true)
     */
    bind(...routes) {
      return bind(this, routes);
    }

    /**@override*/
    copy() {
      return this;
    }

    /**@override*/
    call(port, ...params) {
      const c = call(this, this.#sub(port), ...params);
      c.inner.single = false;
      return c
    }

    #sub(port) {
      if (hasParams(this.query)) throw `cannot nest connection for complete query ${chunks}`
      return port
    }

    /**
     * Adds response mapper to this connection.
     *
     * @param {(any) => any} f - mapping function
     * @returns RpcConn - new instance
     */
    map(f) {
      if (this.mapper) {
        const map = this.mapper;
        this.mapper = arg => f(map(arg));
        return this
      }
      this.mapper = f;
      return this
    }

    /**
     * Writes given method with params to this connection.
     *
     * @async
     * @param {string} method
     * @param {...any} params
     * @returns {Promise<RpcConn>}
     */
    async conn(method, ...params) {
      let cmd = method ? method : "";
      if (params.length > 0) {
        if (cmd) cmd += '?';
        cmd += formatQueryParams(params);
      }
      if (cmd) await this.writeLn(cmd);
      return this
    }

    /**
     * Encodes data into JSON string and writes as line.
     *
     * @async
     * @param {any} data
     * @returns {Promise<undefined>}
     * @throws {string} - IO error message
     */
    async encode(data) {
      let json = JSON.stringify(data);
      if (json === undefined) json = '{}';
      return await super.writeLn(json)
    }

    /**
     * Reads line and parses it into object.
     *
     * @async
     * @returns {Promise<object>} - parsed JSON object
     * @throws {string} - parsed error message
     */
    async decode() {
      const resp = await this.readLn();
      const parsed = JSON.parse(resp);
      if (parsed === null) return null
      if (parsed.error) throw parsed.error
      return parsed
    }

    /**
     * Returns first successfully decoded value or null.
     *
     * @async
     * @returns {Promise<any|null>}
     */
    async request() {
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
     * Collects all decoded values until first null occurrence.
     *
     * @returns {Promise<any[]>} - collected values
     */
    async collect() {
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

  /**
   * Registers given API on {@link RpcClient}.
   *
   * @async
   * @param {RpcClient} client
   * @param {object} api
   * @returns {Promise<void>}
   */
  async function serve(client, api) {
    const listener = await client.register();
    listen(api, listener).catch(bindings.log);
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
      let conn = await listener.accept();
      conn = new RpcConn(conn);
      handle(api, conn).catch(bindings.log).finally(() =>
        conn.close().catch(bindings.log));
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
    const inject = {...api.handlers, ...api.inject, conn: conn};
    const query = conn.query;
    let [handlers, params] = unfold(api.handlers, query);
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
          result = await invoke(handle, inject, params);
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
    const type = typeof handle;
    switch (type) {
      case "function":
        if (!params) return await handle({$: api})
        const [opts, args] = preparePayload(api, params);
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
    const opts = parseQueryParams(params);
    const args = opts._ ? opts._ : [];
    delete opts._;
    opts.$ = api;
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

  /**
   * Adds RPC implementation to {@link ApphostClient}. Compatible with astral query.
   *
   * @extends {ApphostClient}
   */
  class RpcClient extends ApphostClient {

    /**
     * Returns a copy with assigned data.
     *
     * @param data {object}
     * @return RpcClient
     */
    copy(data) {
      return Object.assign(new RpcClient(), {...this, ...data});
    }

    /**
     * Returns a copy with given target id.
     *
     * @param {string} id or alias
     * @return RpcClient
     */
    target(id) {
      return this.copy({targetId: id})
    }

    /**
     * Returns a copy with each route bound as a method.
     *
     * @param {...(string|object)} routes
     * @return {RpcClient&object}
     * @example
     * // A copy of rpc object with target id where foo method is bound to the "foo" path and fiz method is bound to the "bar.baz.fiz" path.
     * const api = rpc.target(id).bind("foo", {"bar.baz": ["fiz"]})
     * await api.foo()
     * await api.fiz("yolo", 1, true)
     */
    bind(...routes) {
      return bind(this, routes);
    }

    /**
     * Creates a call object from this {@link RpcClient}.
     *
     * @param {string} port
     * @param {...arg} params
     * @returns {function&RpcCall}
     */
    call(port, ...params) {
      port = port ? port : "";
      return call(this, port, ...params);
    }

    /**
     * Opens new {@link RpcConn}.
     *
     * @async
     * @param {string} port
     * @param {...arg} params
     * @returns {RpcConn}
     */
    async conn(port, ...params) {
      port = port ? port : "";
      const query = formatQuery(port, params);
      const conn = await super.query(this.targetId, query);
      return new RpcConn(conn)
    }

    /**
     * Serves given API via apphost.
     *
     * @async
     * @param {object} api
     * @returns {Promise<void>}
     */
    async serve(api) {
      return await serve(this, api);
    }
  }

  /**
   * Logs given data to the native console.
   */
  const log = async any => await bindings.log(typeof any == 'object' ? JSON.stringify(any) : any);


  const {
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

  } = bindings;

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
  const apphost = new ApphostClient();

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
  const rpc = new RpcClient();

  exports.apphost = apphost;
  exports.exit = exit;
  exports.log = log;
  exports.platform = platform;
  exports.rpc = rpc;
  exports.sleep = sleep;

  return exports;

})({});
