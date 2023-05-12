const api = new AppHostClient()

log("backend start")

listenHello().catch(err => {
    log("port hello err " + err)
})

async function listenHello() {
    port = await api.listen("hello")
    log("new port " + port.port)
    for (;;) {
        let conn = await port.accept()
        log("new conn " + conn.conn)
        handle(conn).catch(err => {
            log("conn err " + err)
        })
    }
}

async function handle(conn) {
    let data = await conn.read()
    log(data)
    await conn.write("Hello I am backend")
}