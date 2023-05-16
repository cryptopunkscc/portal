class AppHostClient {
    async register(service) {
        await astral_service_register(service)
        return new AppHostListener(service)
    }

    async query(node, query) {
        const conn = await astral_query(node, query)
        return new AppHostConn(conn, query)
    }

    async queryName(node, query) {
        const conn = await astral_query_name(node, query)
        return new AppHostConn(conn, query)
    }

    async nodeInfo(id) {
        return await astral_node_info(id)
    }

    async resolve(name) {
        return await astral_resolve(name)
    }
}

class AppHostListener {
    constructor(port) {
        this.port = port
    }

    async accept() {
        const conn = await astral_conn_accept(this.port)
        return new AppHostConn(conn, this.port)
    }

    async close() {
        await astral_service_close(this.port)
    }
}

class AppHostConn {
    constructor(conn, port) {
        this.conn = conn
        this.port = port
    }

    async read() {
        return await astral_conn_read(this.conn)
    }

    async write(data) {
        return await astral_conn_write(this.conn, data)
    }

    async close() {
        await astral_conn_close(this.conn)
    }
}

const appHost = new AppHostClient()

// ================== RPC extensions ==================

// Bind RPC api of service associated to this connection
AppHostConn.prototype.bindRpc = async function () {
    await astral_rpc_bind_api(this)
}

async function astral_rpc_bind_api(conn) {
    await conn.write(JSON.stringify(["api"]))
    const api = await conn.read()
    log(conn.port + " " + conn.conn + ": == " + api)
    const methods = JSON.parse(api)
    for (let method of methods) {
        conn[method] = async (...data) => {
            const cmd = JSON.stringify([method, ...data])
            log(conn.port + " " + conn.conn + ": => " + cmd)
            await conn.write(cmd)
            const resp = await conn.read()
            const json = JSON.parse(resp)
            log(conn.port + " " + conn.conn + ": <= " + JSON.stringify([method, json]))
            return json
        }
    }
}

// Bind RPC service to given name
AppHostClient.prototype.bindRpc = async function (service, name) {
    await astral_rpc_bind_srv.call(this, service, name)
}

async function astral_rpc_bind_srv(Service, name) {
    const props = Object.getOwnPropertyNames(Service.prototype)
    if (props[0] !== "constructor") throw new Error("Service must have a constructor")
    const methods = props.slice(1, props.length)
    methods.push("api")
    Service.prototype.api = async () => {
        return methods
    }
    const listener = await this.register(name)
    log("listen " + name)
    astral_rpc_listen.call(new Service(), listener).then(log)
}

async function astral_rpc_listen(listener) {
    for (; ;) {
        const conn = await listener.accept()
        log(conn.port + " " + conn.conn + ": accepted")
        astral_rpc_handle.call(this, conn)
    }
}

async function astral_rpc_handle(conn) {
    try {
        for (; ;) {
            const str = await conn.read()
            log(conn.port + " " + conn.conn + ": " + str)
            const query = JSON.parse(str)
            const method = query[0]
            const args = query.slice(1)
            const result = await this[method](...args)
            await conn.write(JSON.stringify(result))
        }
    } catch (e) {
        log(conn.port + " " + conn.conn + ": " + e)
    }
}
