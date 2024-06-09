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

    async query(identity, query) {
      const json = await bindings.astral_query(identity, query);
      const data = JSON.parse(json);
      return new AppHostConn(data, query)
    }

    async queryName(name, query) {
      const json = await bindings.astral_query_name(name, query);
      const data = JSON.parse(json);
      return new AppHostConn(data, query)
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
      return new AppHostConn(data)
    }

    async close() {
      await bindings.astral_service_close(this.port);
    }
  }

  class AppHostConn {
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

  const {log: log$1} = bindings;

  // ================== RPC extensions ==================

  AppHostConn.prototype.jrpcCall = async function (method, ...data) {
    let cmd = method;
    if (data.length > 0) {
      cmd += "?" + JSON.stringify(data);
    }
    // log(this.id + " conn => " + this.query + "." + cmd)
    await this.write(cmd + '\n');
  };

  AppHostConn.prototype.readJson = async function (method) {
    const resp = await this.read();
    const json = JSON.parse(resp);
    return json
  };

  AppHostConn.prototype.rpcQuery = function (method) {
    const conn = this;
    return async function (...data) {
      // log("conn rpc query", method)
      await conn.jrpcCall(method, ...data);
      return await conn.readJson(method)
    }
  };

  AppHostConn.prototype.writeJson = async function (data) {
    const json = JSON.stringify(data);
    // log(this.id + " conn => " + this.query + ":" + json.trimEnd())
    await this.write(json + '\n');
  };

  AppHostConn.prototype.jsonReader = async function (method) {
    const read = async () => await this.readJson(method);
    read.cancel = async () => await this.close();
    return read
  };

  // Bind RPC api of service associated to this connection
  AppHostConn.prototype.bindRpc = async function () {
    await astral_rpc_conn_bind_api(this);
  };

  async function astral_rpc_conn_bind_api(conn) {
    // request api methods
    await conn.jrpcCall("api");

    // read api methods
    const methods = await conn.readJson("api");

    // bind methods
    for (let method of methods) {
      conn[method] = async (...data) => {
        await conn.jrpcCall(method, ...data);
        return await conn.readJson(method)
      };
    }

    // bind subscribe
    conn.subscribe = async (method, ...data) => {
      await conn.jrpcCall(method, ...data);
      return conn.jsonReader(method)
    };
  }

  ApphostClient.prototype.jrpcCall = async function (identity, service, method, ...data) {
    let cmd = service;
    if (method) {
      cmd += "." + method;
    }
    if (data.length > 0) {
      cmd += "?" + JSON.stringify(data);
    }
    const conn = await this.query(identity, cmd);
    // log(conn.id + " client => " + cmd)
    return conn
  };

  ApphostClient.prototype.bindRpc = async function (identity, service) {
    await astral_rpc_client_bind_api(this, identity, service);
    return this
  };

  ApphostClient.prototype.rpcQuery = function (identity, port) {
    const client = this;
    return async function (...data) {
      const conn = await client.jrpcCall(identity, port, "", ...data);
      const json = await conn.readJson(port);
      conn.close().catch(log$1);
      return json
    }
  };

  async function astral_rpc_client_bind_api(client, identity, service) {
    // request api methods
    const conn = await client.jrpcCall(identity, service, "api");

    // read api methods
    const methods = await conn.readJson("api");
    conn.close().catch(log$1);

    // bind methods
    for (let method of methods) {
      client[method] = async (...data) => {
        const conn = await client.jrpcCall(identity, service, method, ...data);
        const json = await conn.readJson(method);
        conn.close().catch(log$1);
        return json
      };
    }

    // bind subscribe
    client.subscribe = async (method, ...data) => {
      const conn = await client.jrpcCall(identity, service, method, ...data);
      return await conn.jsonReader(method)
    };
  }

  // Bind RPC service to given name
  ApphostClient.prototype.bindRpcService = async function (service) {
    return await astral_rpc_bind_srv.call(this, service)
  };

  async function astral_rpc_bind_srv(Service) {
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
    astral_rpc_listen.call(srv, listener).catch(log$1);
    return listener
  }

  async function astral_rpc_listen(listener) {
    for (; ;) {
      const conn = await listener.accept();
      // log(conn.id + " service <= " + conn.query)
      astral_rpc_handle.call(this, conn).catch(log$1);
    }
  }

  async function astral_rpc_handle(conn) {
    try {
      let query = conn.query.slice(this.name.length);
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
          result = await this[method](...args, write, read);
        } catch (e) {
          result = {error: e};
        }
        if (result !== undefined) {
          await conn.writeJson(result);
        }
        if (single) {
          conn.close().catch(log$1);
          break
        }
      }
    } catch (e) {
      // log(conn.id + " service !! " + conn.query + ":" + e)
      conn.close().catch(log$1);
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

  const {log, sleep, platform} = bindings;
  const apphost = new ApphostClient();

  exports.apphost = apphost;
  exports.log = log;
  exports.platform = platform;
  exports.sleep = sleep;

  return exports;

})({});
