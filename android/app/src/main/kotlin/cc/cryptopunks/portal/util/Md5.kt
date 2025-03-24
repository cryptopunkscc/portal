package cc.cryptopunks.portal.util

import java.io.File
import java.io.InputStream
import java.security.MessageDigest

fun File.md5(): String = inputStream().md5()

fun String.md5(): String = byteInputStream().md5()

fun InputStream.md5(): String {
    val md = MessageDigest.getInstance("MD5")
    return use { fis ->
        val buffer = ByteArray(8192)
        generateSequence {
            when (val bytesRead = fis.read(buffer)) {
                -1 -> null
                else -> bytesRead
            }
        }.forEach { bytesRead -> md.update(buffer, 0, bytesRead) }
        md.digest().toHexString(separator = "")
    }
}
