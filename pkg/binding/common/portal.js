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

  class AppHostClient {
    async register(service) {
      await bindings.astral_service_register(service);
      return new AppHostListener(service)
    }

    async query(node, query) {
      const json = await bindings.astral_query(node, query);
      const data = JSON.parse(json);
      return new AppHostConn(data, query)
    }

    async queryName(node, query) {
      const json = await bindings.astral_query_name(node, query);
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
    const cmd = JSON.stringify([method, ...data]);
    log$1(this.query + " " + this.id + ": => " + cmd);
    await this.write(cmd);
  };

  AppHostConn.prototype.readJson = async function (method) {
    const resp = await this.read();
    const json = JSON.parse(resp);
    if (method !== undefined) {
      log$1(this.query + " " + this.id + ": <= " + JSON.stringify([method, json]));
    }
    return json
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

  AppHostClient.prototype.jrpcCall = async function (node, service, method, ...data) {
    const cmd = JSON.stringify([method, ...data]);
    const conn = await this.query(node, service + cmd);
    log$1(service + " " + conn.id + ": => " + cmd);
    return conn
  };

  AppHostClient.prototype.bindRpc = async function (node, service) {
    await astral_rpc_client_bind_api(this, node, service);
    return this
  };

  async function astral_rpc_client_bind_api(client, node, service) {
    // request api methods
    const conn = await client.jrpcCall(node, service, "api");

    // read api methods
    const methods = await conn.readJson("api");
    conn.close().catch(log$1);

    // bind methods
    for (let method of methods) {
      client[method] = async (...data) => {
        const conn = await client.jrpcCall(node, service, method, ...data);
        const json = await conn.readJson(method);
        conn.close().catch(log$1);
        return json
      };
    }

    // bind subscribe
    client.subscribe = async (method, ...data) => {
      const conn = await client.jrpcCall(node, service, method, ...data);
      return await conn.jsonReader(method)
    };
  }

  // Bind RPC service to given name
  AppHostClient.prototype.bindRpcService = async function (service) {

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
    log$1("listen " + srv.name);
    astral_rpc_listen.call(srv, listener).catch(log$1);
    return listener
  }

  async function astral_rpc_listen(listener) {
    for (; ;) {
      const conn = await listener.accept();
      log$1(conn.query + " " + conn.id + ": accepted");
      astral_rpc_handle.call(this, conn).catch(log$1);
    }
  }

  async function astral_rpc_handle(conn) {
    try {
      const send = async (result) =>
        await conn.write(JSON.stringify(result));

      let str = conn.query.slice(this.name.length);
      const single = str.length > 0;
      for (; ;) {
        if (!single) {
          str = await conn.read();
        }
        log$1(this.name + " " + conn.id + ": " + str);
        const query = JSON.parse(str);
        const method = query[0];
        const args = query.slice(1);
        const result = await this[method](...args, send);
        if (result !== undefined) {
          await conn.write(JSON.stringify(result));
        }
        if (single) {
          conn.close().catch(log$1);
          break
        }
      }
    } catch (e) {
      log$1(conn.query + " " + conn.id + ": " + e);
    }
  }

  const {log, sleep, platform} = bindings;
  const apphost = new AppHostClient();

  exports.apphost = apphost;
  exports.log = log;
  exports.platform = platform;
  exports.sleep = sleep;

  return exports;

})({});
