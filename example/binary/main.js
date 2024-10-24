const {apphost, log, sleep} = portal

async function server() {
  const listener = await apphost.register()
  const conn = await listener.accept()
  const buffer = new ArrayBuffer(16)
  const view = new Uint8Array(buffer)
  let n = 0
  for (;;) {
    await sleep(500)
    for (let i = 0; i < buffer.byteLength; i++) {
        view[i] = n++
    }
    log(`server write: ${view}`)
    await conn.write(buffer)
  }
}

async function client() {
  const buffer = new ArrayBuffer(3)
  const view = new Uint8Array(buffer)
  const conn = await apphost.query("example.binary")
  let l = 0
  do {
    await sleep(500)
    l = await conn.read(buffer)
    let actual = []
    for (let i = 0; i < l; i++) {
      actual = [...actual, view[i]]
    }
    log(`client read: ${actual}`)
  } while(l > 0)
}

server().catch(log)
client().catch(log)