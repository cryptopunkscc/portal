package cc.cryptopunks.astral.apphost

import kotlinx.coroutines.CompletableDeferred
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.launch
import java.io.EOFException
import java.io.InputStream
import java.io.OutputStream

suspend inline fun <C : Conn, R> C.use(crossinline block: suspend C.() -> R): R {
    val deferred = CompletableDeferred<R>()
    return coroutineScope {
        launch(Dispatchers.IO) {
            try {
                deferred.complete(block())
            } catch (e: Throwable) {
                deferred.completeExceptionally(e)
            } finally {
                close()
            }
        }
        deferred.await()
    }
}

// =========================== Read ===========================

fun Conn.read(
    size: Int,
) = ByteArray(size)
    .also { buff ->
        val len = read(buff)
        if (len == -1) throw EOFException("EOF")
        check(len == size) { "Expected $size bytes but was $len" }
    }

fun Conn.readMessage(): String? {
    val result = StringBuilder()
    val buffer = ByteArray(4096)
    var len: Int
    do {
        len = read(buffer)
        if (len > 0) result.append(String(buffer.copyOf(len)))
    } while (len == buffer.size)
    return when {
        len == -1 && result.isEmpty() -> null
        else -> result.toString().takeIf { it != "null" }
    }
}

// =========================== Write ===========================

fun Conn.write(input: InputStream) {
    val size = 16 * 1024
    val buffer = ByteArray(size)
    var len: Int
    while (true) {
        len = input.read(buffer)
        when (len) {
            size -> write(buffer)
            -1 -> break
            else -> write(buffer.copyOf(len))
        }
    }
}

fun Conn.write(
    bytes: ByteArray,
    size: ByteArray,
) {
    write(size)
    write(bytes)
}

fun Conn.write(
    bytes: ByteArray,
    formatSize: Int.() -> ByteArray,
) = write(
    bytes = bytes.size.formatSize(),
    size = bytes
)

var Conn.byte: Byte
    get() = read(Byte.SIZE_BYTES).byte
    set(value) {
        write(value.bytes)
    }

var Conn.short: Short
    get() = read(Short.SIZE_BYTES).short
    set(value) {
        write(value.bytes)
    }

var Conn.int: Int
    get() = read(Int.SIZE_BYTES).int
    set(value) {
        write(value.bytes)
    }

var Conn.long: Long
    get() = read(Long.SIZE_BYTES).long
    set(value) {
        write(value.bytes)
    }

var Conn.bytes8: ByteArray
    get() = read(read(Byte.SIZE_BYTES).byte.toUByte().toInt())
    set(bytes) {
        write(bytes, bytes.size.toByte().bytes)
    }

var Conn.bytes16: ByteArray
    get() = read(read(Short.SIZE_BYTES).short.toUShort().toInt())
    set(bytes) {
        write(bytes, bytes.size.toShort().bytes)
    }

var Conn.bytes32: ByteArray
    get() = read(read(Int.SIZE_BYTES).int)
    set(bytes) {
        write(bytes, bytes.size.bytes)
    }

var Conn.string8: String
    get() = bytes8.decodeToString()
    set(string) {
        bytes8 = string.encodeToByteArray()
    }

var Conn.string16: String
    get() = bytes16.decodeToString()
    set(string) {
        bytes16 = string.encodeToByteArray()
    }

var Conn.string32: String
    get() = bytes32.decodeToString()
    set(string) {
        bytes32 = string.encodeToByteArray()
    }

var Conn.identity: ByteArray
    get() {
        val id = read(33)
        return when {
            id.all { it == zero } -> ByteArray(0)
            else -> id
        }
    }
    set(value) {
        val id = when {
            value.isEmpty() || value.contentEquals(localnode) -> ByteArray(33)
            value.size == 33 -> value
            else -> throw IllegalArgumentException("Invalid identity size ${value.size}, have to be 33 or 0")
        }
        write(id)
    }

private val zero = 0.toByte()
private val localnode = "localnode".toByteArray()

// =========================== Stream ===========================

val Conn.inputStream: InputStream get() = ConnInputStream(this)
val Conn.outputStream: OutputStream get() = ConnOutputStream(this)

private class ConnInputStream(private val conn: Conn) : InputStream() {
    override fun readNBytes(len: Int) = conn.readN(len)
    override fun read(): Int = conn.readN(1)[0].toInt()
    override fun read(buff: ByteArray) = conn.read(buff)
    override fun read(b: ByteArray, off: Int, len: Int): Int {
        if (off == 0 && len == b.size) {
            return conn.read(b)
        }
        val r = conn.readN(len)
        System.arraycopy(r, 0, b, off, r.size)
        return r.size
    }
}

private class ConnOutputStream(private val conn: Conn) : OutputStream() {

    override fun write(p0: Int) {
        conn.write(ByteArray(1) { p0.toByte() })
    }

    override fun write(b: ByteArray, off: Int, len: Int) {
        if (off == 0 && len == b.size) {
            conn.write(b)
            return
        }
        val copy = b.copyOfRange(off, off + len)
        conn.write(copy)
    }

    override fun close() {
        conn.close()
    }
}
