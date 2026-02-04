package cc.cryptopunks.astral.apphost

interface AppHostClient {
    fun query(query: String, nodeId: String = ""): Conn
    fun register(name: String): ApphostListener
    fun resolve(name: String): String
}

interface ApphostListener {
    operator fun next(): QueryData
}

interface QueryData {
    fun accept(): Conn
    fun caller(): String
    fun query(): String
    fun reject()
}

interface Conn {
    fun close()
    fun read(buff: ByteArray): Int
    fun readN(n: Int): ByteArray
    fun write(buff: ByteArray): Int
}

abstract class AppHostClientAdapter<C : Conn>(
    private val client: AppHostClient,
) : AppHostClient by client {

    protected abstract fun convert(conn: Conn): C

    override fun query(query: String, nodeId: String): C = convert(client.query(query, nodeId))

    override fun register(name: String): ApphostListener = Listener(client.register(name))

    inner class Listener(private val listener: ApphostListener) : ApphostListener {
        override fun next(): QueryData = Query(listener.next())
    }

    inner class Query(private val query: QueryData) : QueryData by query {
        override fun accept() = convert(query.accept())
    }
}
