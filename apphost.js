class AppHostClient {
    async register(service) {
        await astral_service_register(service)
        return new AppHostListener(service)
    }

    async query(node, query) {
        const conn = await astral_query(node, query)
        return new AppHostConn(conn)
    }

    async queryName(node, query) {
        const conn = await astral_query_name(node, query)
        return new AppHostConn(conn)
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
        return new AppHostConn(conn)
    }

    async close() {
        await astral_service_close(this.port)
    }
}

class AppHostConn {
    constructor(conn) {
        this.conn = conn
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