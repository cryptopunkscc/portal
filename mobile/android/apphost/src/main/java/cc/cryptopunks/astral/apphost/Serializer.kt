package cc.cryptopunks.astral.apphost

interface Serializer {
    fun encode(any: Any): String
    fun <T> decode(string: String, type: Class<T>): T
    fun <T> decodeList(string: String, type: Class<T>): List<T>
    fun <K, V> decodeMap(string: String, key: Class<K>, value: Class<V>): Map<K, V>
}

fun Conn.serializer(serializer: Serializer) = ConnSerializer(this, serializer)

data class ConnSerializer(
    val conn: Conn,
    val serializer: Serializer,
) : Serializer by serializer,
    Conn by conn

fun ConnSerializer.encodeL8(any: Any) {
    string8 = encode(any)
}

fun ConnSerializer.encodeL16(any: Any) {
    string16 = encode(any)
}

fun ConnSerializer.encodeL32(any: Any) {
    string32 = encode(any)
}

fun ConnSerializer.encodeLine(any: Any) {
    write((encode(any) + "\n").encodeToByteArray())
}

fun <T> ConnSerializer.decode8(type: Class<T>): T =
    decode(string8, type)

fun <T> ConnSerializer.decode16(type: Class<T>): T =
    decode(string16, type)

fun <T> ConnSerializer.decode32(type: Class<T>): T =
    decode(string32, type)

fun <T> ConnSerializer.decodeMessage(type: Class<T>): T =
    decode(readMessage().orEmpty(), type)

inline fun <reified T> ConnSerializer.decodeMessage(): T =
    decode(readMessage().orEmpty(), T::class.java)

inline fun <reified T> ConnSerializer.decodeList(): List<T> =
    decodeList(readMessage().orEmpty(), T::class.java)

inline fun <reified K, reified V> ConnSerializer.decodeMap(): Map<K, V> =
    decodeMap(readMessage().orEmpty(), K::class.java, V::class.java)

inline fun <reified T> ConnSerializer.decode8(): T =
    decode(string8, T::class.java)

inline fun <reified T> ConnSerializer.decode16(): T =
    decode(string16, T::class.java)

inline fun <reified T> ConnSerializer.decode32(): T =
    decode(string32, T::class.java)

inline fun <reified T> ConnSerializer.decodeList8(): List<T> =
    decodeList(string8, T::class.java)

inline fun <reified T> ConnSerializer.decodeList16(): List<T> =
    decodeList(string16, T::class.java)

inline fun <reified T> ConnSerializer.decodeList32(): List<T> =
    decodeList(string32, T::class.java)
