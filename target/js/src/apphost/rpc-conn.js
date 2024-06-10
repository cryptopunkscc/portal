import {AppHostConn} from "./adapter.js";
import {bindings} from "../bindings.js";

const {log} = bindings


AppHostConn.prototype.jrpcCall = async function (method, ...data) {
  let cmd = method
  if (data.length > 0) {
    cmd += "?" + JSON.stringify(data)
  }
  // log(this.id + " conn => " + this.query + "." + cmd)
  await this.write(cmd + '\n')
}

AppHostConn.prototype.readJson = async function (method) {
  const resp = await this.read()
  const json = JSON.parse(resp)
  if (method !== undefined) {
    // log(this.id + " conn <= " + this.query  + ":" + resp.trimEnd())
  }
  return json
}

AppHostConn.prototype.rpcQuery = function (method) {
  const conn = this
  return async function (...data) {
    // log("conn rpc query", method)
    await conn.jrpcCall(method, ...data)
    return await conn.readJson(method)
  }
}

AppHostConn.prototype.writeJson = async function (data) {
  const json = JSON.stringify(data)
  // log(this.id + " conn => " + this.query + ":" + json.trimEnd())
  await this.write(json + '\n')
}

AppHostConn.prototype.jsonReader = async function (method) {
  const read = async () => await this.readJson(method)
  read.cancel = async () => await this.close()
  return read
}

// Bind RPC api of service associated to this connection
AppHostConn.prototype.bindRpc = async function () {
  const conn = this
  // request api methods
  await conn.jrpcCall("api")

  // read api methods
  const methods = await conn.readJson("api")

  // bind methods
  for (let method of methods) {
    conn[method] = async (...data) => {
      await conn.jrpcCall(method, ...data)
      return await conn.readJson(method)
    }
  }

  // bind subscribe
  conn.subscribe = async (method, ...data) => {
    await conn.jrpcCall(method, ...data)
    return conn.jsonReader(method)
  }
}