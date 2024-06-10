import {ApphostClient} from "./adapter.js";
import {bindings} from "../bindings.js";

const {log} = bindings


// Bind RPC service to given name
ApphostClient.prototype.bindRpcService = async function (Service) {
  const props = Object.getOwnPropertyNames(Service.prototype)
  if (props[0] !== "constructor") throw new Error("Service must have a constructor")
  const methods = props.slice(1, props.length)
  methods.push("api")
  Service.prototype.api = async () => {
    return methods
  }
  const srv = new Service()
  const listener = await this.register(srv.name + "*")
  // log("listen " + srv.name)
  astral_rpc_listen(srv, listener).catch(log)
  return listener
}

async function astral_rpc_listen(srv, listener) {
  for (; ;) {
    const conn = await listener.accept()
    // log(conn.id + " service <= " + conn.query)
    astral_rpc_handle(srv, conn).catch(log)
  }
}

async function astral_rpc_handle(srv, conn) {
  try {
    let query = conn.query.slice(srv.name.length)
    let method = query, args = []
    const single = query !== ''
    const write = async (data) => await conn.writeJson(data)
    const read = async (method) => await conn.readJson(method)

    for (; ;) {
      if (!single) {
        query = await conn.read();
        // log(conn.id + " service <== " + query)
      }
      [method, args] = parseQuery(query)

      let result
      try {
        result = await srv[method](...args, write, read)
      } catch (e) {
        result = {error: e}
      }
      if (result !== undefined) {
        await conn.writeJson(result)
      }
      if (single) {
        conn.close().catch(log)
        break
      }
    }
  } catch (e) {
    // log(conn.id + " service !! " + conn.query + ":" + e)
    conn.close().catch(log)
  }
}

function parseQuery(query) {
  if (query[0] === '.') {
    query = query.slice(1)
  }
  const match = /[?\[{]/.exec(query)
  const method = query.slice(0, match.index)
  let payload = query.slice(match.index)
  if (payload[0] === '?') {
    payload = payload.slice(1)
  }
  let args = []
  if (payload) {
    args = JSON.parse(payload)
  }
  return [method, args]
}
