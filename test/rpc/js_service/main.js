const apphost = portal.apphost
const log = portal.log

class Service {

  constructor() {
    this.name = "flow"
  }

  async func1(msg, fail){
    if (fail) {
      throw msg
    }
    return msg
  }

  async func2(b, i, f, s){
    return [b, i, f, s]
  }

  async func3(struct){
    return struct
  }

  async func4(b, i, f, s){
    return {b: b, i: i, f: f, s: s}
  }
}

apphost.bindRpcService(Service).catch(log)
