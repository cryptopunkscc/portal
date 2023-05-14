const api = new AppHostClient()

log("backend start")

api.listen("hello")
    .catch(log)
    .then(handle)

function handle(port) {
    log("new port " + port.port)
    port.accept().catch(log).then(conn => {
        log("new conn " + conn.conn)
        conn.read().catch(log).then(data => {
            log(data)
            conn.write("Hello I am backend").catch(log)
            handle(port)
        })
        // TODO async handle doesn't work for some reason. Investigation needed.
        // handle(port)
    })
}
