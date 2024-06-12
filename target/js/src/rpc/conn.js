import {ApphostConn} from "../apphost/adapter.js";

export class RpcConn extends ApphostConn {
  constructor(data) {
    super(data);
  }

  async encode(data) {
    let json = JSON.stringify(data)
    return await super.write(json + '\n')
  }

  async decode() {
    const resp = await this.read()
    return JSON.parse(resp)
  }

  async call(method, ...params) {
    let cmd = method
    if (params) {
      cmd += '?' + JSON.stringify(params)
    }
    await this.write(cmd + '\n')
  }

  async request(query, ...params) {
    await this.call(query, ...params)
    return await this.decode()
  }

  caller(method) {
    return async (...params) => await this.call(method, ...params)
  }

  requester(method) {
    return async (...params) => await this.request(method, ...params)
  }

  bind(methods) {
    this.boundMethods = methods
    for (let method of methods) {
      this[method] = this.requester(method)
    }
    return this
  }
}
