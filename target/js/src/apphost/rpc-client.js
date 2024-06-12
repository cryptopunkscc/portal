import {ApphostClient} from "./adapter.js";
import {bindings} from "../bindings.js";
import "./rpc-conn.js"

const {log} = bindings

ApphostClient.prototype.rpcCall = async function (identity, service, method, ...data) {
  let cmd = service
  if (method) {
    cmd += "." + method
  }
  if (data.length > 0) {
    cmd += "?" + JSON.stringify(data)
  }
  const conn = await this.query(cmd, identity)
  // log(conn.id + " client => " + cmd)
  return conn
}

ApphostClient.prototype.rpcQuery = function (identity, port) {
  const client = this
  return async function (...data) {
    const conn = await client.rpcCall(identity, port, "", ...data)
    const json = await conn.readJson(port)
    conn.close().catch(log)
    return json
  }
}

ApphostClient.prototype.bindRpc = async function (identity, service) {
  const client = this
  // request api methods
  const conn = await client.rpcCall(identity, service, "api")

  // read api methods
  const methods = await conn.readJson("api")
  conn.close().catch(log)

  // bind methods
  for (let method of methods) {
    client[method] = async (...data) => {
      const conn = await client.rpcCall(identity, service, method, ...data)
      const json = await conn.readJson(method)
      conn.close().catch(log)
      return json
    }
  }

  // bind subscribe
  client.subscribe = async (method, ...data) => {
    const conn = await client.rpcCall(identity, service, method, ...data)
    return await conn.jsonReader(method)
  }
  return client
}
