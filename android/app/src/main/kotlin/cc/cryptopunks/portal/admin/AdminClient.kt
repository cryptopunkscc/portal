package cc.cryptopunks.portal.admin

import cc.cryptopunks.astral.bind.astral.JsAppHostClient
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.isActive

class AdminClient(
    private val client: JsAppHostClient,
) {
    private var id = ""

    private val accumulated = MutableStateFlow("")

    val output = accumulated.asStateFlow()

    suspend fun connect() = coroutineScope {
        while (isActive) try {
            id = client.query("", "admin")
            if (accumulated.value.isBlank())
                client.connWrite(id, "help\n")
            while (isActive) {
                val string = client.connRead(id).clearFormatting()
                accumulated.update { it + string }
            }
        } catch (e: Exception) {
            client.runCatching { connClose(id) }
            val message = e.message.orEmpty()
            when {
                "apphost" in message -> break
                "EOF" in message -> break
                else -> throw e
            }
        }
    }

    fun query(text: String) {
        accumulated.update { it + text + "\n" }
        client.connWrite(id, text + "\n")
    }
}

private fun String.clearFormatting() = replace(formattingRegex, "")

private val formattingRegex = Regex("""\[\d+m|""")
