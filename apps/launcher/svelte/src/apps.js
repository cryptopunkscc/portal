import {log} from "portal";
import portal from "./portal.js";
import {writable} from "svelte/store";

export class AppsRepository {
  constructor() {
    this.apps = []
    this.store = writable([])
    this.channel = null
    portal.observe().catch(console.log).then(channel => {
      this.channel = channel
      this.run().catch(console.log)
      this.loadMore()
    })
    console.log(this)
  }

  async run() {
    for (; this.channel;) {
      let app;
      try {
        app = await this.channel.next()
      } catch (e) {
        log("error: ", JSON.stringify(e))
        throw e
      }
      this.apps.push(app)
      this.store.set(this.apps)
    }
    log("close run")
  }

  subscribe(run, invalidate) {
    // log("subscribe ", this.channel)
    return this.store.subscribe(run, invalidate)
  }

  loadMore(num) {
    // log("loadMore")
    num = num || 10
    this.channel?.more(num)?.catch(console.log)
  }

  cancel() {
    log("cancel")
    this.channel?.close()?.catch(console.log)
    this.channel = null
  }
}
