export class Service {

  constructor() {
    this.name = "rpc"
    this.counter = 0
  }

  async get(arg) {
    return {
      arg: arg,
      val: "Hello RPC"
    }
  }

  async sum(a, b) {
    return a + b
  }

  async inc() {
    return ++this.counter
  }

  async ticker(send) {
    let counter = 100
    for (; ;) {
      await send(counter++)
      await sleep(1000)
    }
  }
}