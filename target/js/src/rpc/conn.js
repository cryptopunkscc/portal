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

  async observe(consume) {
    for (;;) {
      const next = await this.decode()
      const last = await consume(next)
      this.value = next
      if (last) {
        return last
      }
    }
  }

  caller(method) {
    return async (...params) => await this.call(method, ...params)
  }

  requester(method) {
    return async (...params) => await this.request(method, ...params)
  }

  observer(method) {
    return (...params) => ({
      observe: async (consume) => {
        await this.call(method, ...params)
        return await this.observe(consume)
          .finally(() => this.close()) // TODO consider if it's ok to close conn automatically.
      }
    })
  }

  bind(...methods) {
    for (let method of methods) {
      const collect = method[0] === '*'
      if (collect) {
        method = method.split(1)
      }
      if (this[method]) {
        throw `method '${method}' already exist`
      }
      if (collect) {
        this[method] = this.observer(method)
      } else {
        this[method] = this.requester(method)
      }
    }
    return this
  }
}
