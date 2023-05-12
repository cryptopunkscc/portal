class AppHostClient {
    async listen(port) {
        await astral_port_listen(port)
        return new Port(port)
    }

    async dial(node, query) {
        const conn = await astral_dial(node, query)
        return new Conn(conn)
    }

    async dialName(node, query) {
        const conn = await astral_dial_name(node, query)
        return new Conn(conn)
    }

    async nodeInfo(id) {
        return await astral_node_info(id)
    }

    async resolve(name) {
        return await astral_resolve(name)
    }
}

class Port {
    constructor(port) {
        this.port = port
    }

    async accept() {
        const conn = await astral_conn_accept(this.port)
        return new Conn(conn)
    }

    async close() {
        astral_port_close(this.port)
    }
}

class Conn {
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
        astral_conn_close(this.conn)
    }
}