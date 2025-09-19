package cc.cryptopunks.portal.util

import java.io.File
import java.io.InputStream
import java.security.MessageDigest

fun File.md5(): String = inputStream().md5()

fun String.md5(): String = byteInputStream().md5()

fun InputStream.md5(): String = use { inputStream ->
    MessageDigest.getInstance("MD5").apply {
        val buffer = ByteArray(8192)
        while (true) {
            val len = inputStream.read(buffer)
            if (len == -1) break
            update(buffer, 0, len)
        }
    }.digest().toHexString(separator = "")
}
