package cc.cryptopunks.portal.exception

import android.content.Context
import cc.cryptopunks.portal.util.FileDateFormat
import cc.cryptopunks.portal.util.md5
import java.util.Date

class ExceptionStorage(context: Context) {

    private val errorsDir = context.cacheDir.resolve("errors").apply { mkdirs() }

    fun save(e: Throwable) = try {
        val stackTraceString = e.stackTraceToString()
        val datePrefix = FileDateFormat.format(Date())
        val exceptionName = e.message
            ?.run { filter { it.isLetterOrDigit() || it == ' ' }.replace(' ', '_') }
            ?: stackTraceString.md5()
        val file = errorsDir.resolve(datePrefix + exceptionName)
        file.writeText(stackTraceString)
    } catch (e2: Throwable) {
        e2.addSuppressed(e)
        e2.printStackTrace()
    }
}
