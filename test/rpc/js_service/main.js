const apphost = portal.apphost
const log = portal.log

class Service {

  async func1(msg, fail) {
    if (fail) {
      throw msg
    }
    return msg
  }

  async func2(b, i, f, s) {
    return [b, i, f, s]
  }

  async func3(struct) {
    return struct
  }

  async func4(b, i, f, s) {
    return {b: b, i: i, f: f, s: s}
  }
}

const service = new Service();

apphost.registerRpc({
  routes: [
    "flow*",
    "request:",
    "request.func1",
    "request.func2",
    "request.func3",
    "request.func4",
  ],
  handlers: {
    request: service,
    flow: service,
  },
  inject: {}
}).catch(log)

