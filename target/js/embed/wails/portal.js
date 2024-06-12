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

  /* eslint-disable */
  const platform$1 = typeof window['go'] === "undefined" ? undefined : "wails";

  /* eslint-disable */
  const adapter = () => ({
    // apphost
    astral_conn_accept: window['go']['main']['Adapter']['ConnAccept'],
    astral_conn_close: window['go']['main']['Adapter']['ConnClose'],
    astral_conn_read: window['go']['main']['Adapter']['ConnRead'],
    astral_conn_write: window['go']['main']['Adapter']['ConnWrite'],
    astral_node_info: window['go']['main']['Adapter']['NodeInfo'],
    astral_query: window['go']['main']['Adapter']['Query'],
    astral_query_name: window['go']['main']['Adapter']['QueryName'],
    astral_resolve: window['go']['main']['Adapter']['Resolve'],
    astral_service_close: window['go']['main']['Adapter']['ServiceClose'],
    astral_service_register: window['go']['main']['Adapter']['ServiceRegister'],
    astral_interrupt: window['go']['main']['Adapter']['Interrupt'],
    // runtime
    sleep: window['go']['main']['Adapter']['Sleep'],
    log: (...arg) => window['go']['main']['Adapter']['LogArr'](arg),
  });

  inject(platform$1, adapter);

  // ================== Object oriented adapter ==================

  class ApphostClient {
    async register(service) {
      await bindings.astral_service_register(service);
      return new AppHostListener(service)
    }

    async query(query, identity) {
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
      return await bindings.astral_conn_read(this.id)
    }

    async write(data) {
      return await bindings.astral_conn_write(this.id, data)
    }

    async close() {
      await bindings.astral_conn_close(this.id);
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

  ApphostClient.prototype.rpcQuery = function (identity, port) {
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

  exports.apphost = apphost;
  exports.log = log;
  exports.platform = platform;
  exports.sleep = sleep;

  return exports;

})({});
