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
const platform$3 = typeof window['go'] === "undefined" ? undefined : "wails";

/* eslint-disable */
const adapter$2 = () => ({
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
  const conn = this;
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
};

const {log: log$2} = bindings;

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

ApphostClient.prototype.rpcQuery = function (identity, port) {
  const client = this;
  return async function (...data) {
    const conn = await client.jrpcCall(identity, port, "", ...data);
    const json = await conn.readJson(port);
    conn.close().catch(log$2);
    return json
  }
};

ApphostClient.prototype.bindRpc = async function (identity, service) {
  const client = this;
  // request api methods
  const conn = await client.jrpcCall(identity, service, "api");

  // read api methods
  const methods = await conn.readJson("api");
  conn.close().catch(log$2);

  // bind methods
  for (let method of methods) {
    client[method] = async (...data) => {
      const conn = await client.jrpcCall(identity, service, method, ...data);
      const json = await conn.readJson(method);
      conn.close().catch(log$2);
      return json
    };
  }

  // bind subscribe
  client.subscribe = async (method, ...data) => {
    const conn = await client.jrpcCall(identity, service, method, ...data);
    return await conn.jsonReader(method)
  };
  return client
};

const {log: log$1} = bindings;


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
  astral_rpc_listen(srv, listener).catch(log$1);
  return listener
};

async function astral_rpc_listen(srv, listener) {
  for (; ;) {
    const conn = await listener.accept();
    // log(conn.id + " service <= " + conn.query)
    astral_rpc_handle(srv, conn).catch(log$1);
  }
}

async function astral_rpc_handle(srv, conn) {
  try {
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

export { apphost, log, platform, sleep };
