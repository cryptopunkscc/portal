const {apphost, log, sleep} = portal

log("backend start args: " + args)

listen().catch(err => {
  log("listen err: " + err)
})

call().catch(err => {
  log("call err: " + err)
})

async function listen() {
  log("listener registered")
  const listener = await apphost.register()
  const conn = await listener.accept()
  log("conn accepted: " + conn.id)
  const data = await conn.readLn()
  await conn.writeLn(data)
  await conn.close()
  await listener.close()
}

async function call() {
  await sleep(200)
  let conn = await portal.apphost.query("test.basic.js", "portal")
  await conn.writeLn("hello")
  let result = await conn.readLn()
  log("result: " + result)
  await conn.close()
  portal.exit(0)
}



