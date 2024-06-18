const {log, rpc, sleep} = portal

class Service {

  constructor() {
    this.counter = 0
  }

  // Example
  async link(ctx) {
    await ctx.conn.encode()
  }

  // Example observable
  async ticker(limit, delay, {conn}) {

    // Example validation
    const l = Number(limit)
    const d = Number(delay)
    if (isNaN(l)) {
      throw `limit must be a number instead of ${typeof limit}`
    }
    if (l < 0) {
      throw "limit cannot be negative"
    }
    if (isNaN(d)) {
      throw `delay must be a number instead of ${typeof delay}`
    }
    if (d < 0) {
      throw `delay cannot be negative`
    }

    new Promise(() => conn.decode())
    log("start ticker")
    let counter = 0
    while (!conn.done && counter <= limit) {
      await conn.encode(counter++)
      await sleep(delay)
    }
    return null
  }
}

export const serve = async () => rpc.serve({
  routes: ["link", "ticker"],
  handlers: new Service(),
})