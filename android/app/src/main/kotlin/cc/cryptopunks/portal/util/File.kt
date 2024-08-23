package cc.cryptopunks.astral.agent.util

import android.os.Build
import android.os.FileObserver
import java.io.File

fun File.fileObserver(
    mask: Int,
    onEvent: (event: Int, path: String?) -> Unit,
): FileObserver = when {

    Build.VERSION.SDK_INT < Build.VERSION_CODES.Q
    -> object : FileObserver(absolutePath, mask) {
        override fun onEvent(event: Int, path: String?) = onEvent(event, path)
    }

    else
    -> object : FileObserver(this, mask) {
        override fun onEvent(event: Int, path: String?) = onEvent(event, path)
    }
}

internal fun File.createBackup(
    prefix: String = FileDateFormat.format(System.currentTimeMillis()),
) {
    if (!exists()) return
    val dir = parentFile ?: return
    val previous = dir.listFiles()?.dropLastWhile { it == this }?.lastOrNull()
    if (md5() == previous?.md5()) return

    val backupName = prefix + name
    val backup = dir.resolve(backupName)

    copyTo(backup)
}
