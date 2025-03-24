package cc.cryptopunks.portal.util

fun ByteArray.toHexString(separator: String = ":"): String =
    joinToString(separator) { String.format("%02X", it) }