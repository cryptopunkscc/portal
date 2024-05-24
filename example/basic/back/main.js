const {apphost, log, sleep} = portal

log("backend start")

listenHello().catch(err => {
  log("port hello err " + err)
})

async function listenHello() {
  this.port = await apphost.register("")
  log("new port " + this.port.port)
  for (; ;) {
    let conn = await this.port.accept()
    log("new conn " + conn.id)
    handle(conn).catch(err => {
      log("conn err " + err)
    })
  }
}

async function handle(conn) {
  let data = await conn.read()
  log("blocking " + conn.id)
  await sleep(3000)
  log(data)
  await conn.write("Hello I am backend\n")
  await conn.close()
}
