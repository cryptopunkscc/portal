const {sleep, log} = portal

export class Service {

  constructor() {
    this.name = "portal"
    this.counter = 0
  }

  async open(id) {
    log("open " + id)
    return id
  }

  async launch(id) {
    log("launch " + id)
  }

  async install(id) {
    log("install " + id)
  }

  async uninstall(id) {
    log("uninstall " + id)
  }

  async observe(write, read){
    log("======================== observe \n" + write + "\n" + read)
    let end = 0
    let begin = end
    for (;;) {
      log("read1")
      const num = await read()
      log("read2")
      log(num)
      begin = end
      end = end + num
      for (let i = begin; i < end; i++) {
        let next = {id: i, name: "app" + i}
        await write(next)
        log("next: ", next)
      }
    }
  }
}
