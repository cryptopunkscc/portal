const {log, rpc, sleep} = portal

class Service {

  constructor() {
    this.counter = 0
  }

  // Example link for
  async link(ctx) {
    await ctx.conn.encode()
  }

  // Example observable handler
  async ticker(limit, delay, {conn}) {

    // Example validation
    const l = Number(limit)
    const d = Number(delay)
    if (isNaN(l)) throw `limit must be a number instead of ${typeof limit}`
    if (l < 0) throw "limit cannot be negative"
    if (isNaN(d)) throw `delay must be a number instead of ${typeof delay}`
    if (d < 0) throw `delay cannot be negative`

    // Read connection close on separate routine to detect client disconnection ASAP.
    conn.decode()

    // Serve ticket.
    log("start ticker")
    let counter = 0
    while (!conn.done && counter <= limit) {
      await conn.encode(counter++)
      await sleep(delay)
    }
    return null
  }
}

// Export server initializer
export const serve = async () => rpc.serve({
  handlers: new Service(),
  routes: ["link", "ticker"],
})