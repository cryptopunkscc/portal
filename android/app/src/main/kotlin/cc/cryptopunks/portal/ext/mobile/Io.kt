package cc.cryptopunks.portal.ext.mobile

import cc.cryptopunks.portal.core.mobile.AsyncReader
import cc.cryptopunks.portal.core.mobile.AsyncWriter
import cc.cryptopunks.portal.core.mobile.CallbackN
import cc.cryptopunks.portal.core.mobile.Closer
import cc.cryptopunks.portal.core.mobile.Conn
import cc.cryptopunks.portal.core.mobile.ReadCloser
import cc.cryptopunks.portal.core.mobile.Reader
import cc.cryptopunks.portal.core.mobile.Writer
import kotlinx.coroutines.CancellableContinuation
import kotlinx.coroutines.suspendCancellableCoroutine
import java.io.IOException
import java.io.InputStream
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException

fun Conn.coroutine(): CoroutineReadWriteCloser = CoroutineReadWriteCloserImpl(this, this, this)

interface CoroutineReader {
    suspend fun read(buff: ByteArray): Long
}

interface CoroutineWriter {
    suspend fun write(buff: ByteArray): Long
}

interface CoroutineReadWriter : CoroutineReader, CoroutineWriter
interface CoroutineReadCloser : CoroutineReader, Closer
interface CoroutineWriteCloser : CoroutineWriter, Closer
interface CoroutineReadWriteCloser : CoroutineReadWriter, CoroutineReadCloser, CoroutineWriteCloser

class CoroutineReadWriteCloserImpl(
    reader: Reader,
    writer: Writer,
    closer: Closer,
) : CoroutineReadWriteCloser,
    CoroutineReader by CoroutineReaderImpl(reader),
    CoroutineWriter by CoroutineWriterImpl(writer),
    Closer by closer


class CoroutineReaderImpl(reader: Reader) : CoroutineReader {
    private val async = AsyncReader(reader)
    override suspend fun read(buff: ByteArray): Long = suspendCancellableCoroutine {
        async.read(buff, callback(it))
    }
}

class CoroutineWriterImpl(writer: Writer) : CoroutineWriter {
    private val async = AsyncWriter(writer)
    override suspend fun write(buff: ByteArray): Long = suspendCancellableCoroutine {
        async.write(buff, callback(it))
    }
}

fun callback(continuation: CancellableContinuation<Long>) = CallbackN { n, err ->
    if (err == null) continuation.resume(n)
    else continuation.resumeWithException(err)
}

fun ReadCloser.inputStream(): InputStream = ReaderInputStream(this)

private class ReaderInputStream(private val readCloser: ReadCloser) : InputStream() {
    override fun read(): Int = catchEOF(-1) { readCloser.readN(1).first().toInt() }
    override fun read(b: ByteArray): Int = catchEOF(-1) { readCloser.read(b).toInt() }
    override fun readNBytes(len: Int): ByteArray = catchEOF(ByteArray(0)) { readCloser.readN(len.toLong()) }
    override fun readNBytes(b: ByteArray, off: Int, len: Int): Int = catchEOF(0) {
        val b2 = readCloser.readN(len.toLong())
        System.arraycopy(b2, 0, b, off, b2.size)
        b2.size
    }

    override fun close() = catchEOF(Unit) { readCloser.close() }

    private fun <T> catchEOF(default: T, block: () -> T): T = try {
        block()
    } catch (e: Throwable) {
        if (e.message == "EOF") default
        else throw IOException(e)
    }
}