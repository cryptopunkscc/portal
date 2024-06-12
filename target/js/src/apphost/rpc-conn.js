import {ApphostConn} from "./adapter.js";
import {bindings} from "../bindings.js";

const {log} = bindings


ApphostConn.prototype.readJson = async function (method) {
  const resp = await this.read()
  const json = JSON.parse(resp)
  if (method !== undefined) {
    // log(this.id + " conn <= " + this.query  + ":" + resp.trimEnd())
  }
  return json
}

ApphostConn.prototype.jsonReader = async function (method) {
  const read = async () => await this.readJson(method)
  read.cancel = async () => await this.close()
  return read
}

ApphostConn.prototype.writeJson = async function (data) {
  // if (Array.isArray(data) && data.length === 1) {
  //   data = data[0]
  // }
  const json = JSON.stringify(data)
  // log(this.id + " conn => " + this.query + ":" + json.trimEnd())
  await this.write(json + '\n')
}

ApphostConn.prototype.rpcCall = async function (method, ...data) {
  let cmd = method
  if (data.length > 0) {
    cmd += "?" + JSON.stringify(data)
  }
  // log(this.id + " conn => " + this.query + "." + cmd)
  await this.write(cmd + '\n')
}

ApphostConn.prototype.rpcQuery = function (method) {
  const conn = this
  return async function (...data) {
    // log("conn rpc query", method)
    await conn.rpcCall(method, ...data)
    return await conn.readJson(method)
  }
}

// Bind RPC api of service associated to this connection
ApphostConn.prototype.bindRpc = async function () {
  const conn = this
  // request api methods
  await conn.rpcCall("api")

  // read api methods
  const methods = await conn.readJson("api")

  // bind methods
  for (let method of methods) {
    conn[method] = async (...data) => {
      await conn.rpcCall(method, ...data)
      return await conn.readJson(method)
    }
  }

  // bind subscribe
  conn.subscribe = async (method, ...data) => {
    await conn.rpcCall(method, ...data)
    return conn.jsonReader(method)
  }
}