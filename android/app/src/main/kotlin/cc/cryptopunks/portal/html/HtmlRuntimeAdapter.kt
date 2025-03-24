package cc.cryptopunks.portal.html

import android.util.Log
import android.webkit.JavascriptInterface
import android.webkit.WebView
import cc.cryptopunks.portal.core.bind.Core
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.launch
import org.json.JSONObject
import java.util.UUID

internal class HtmlRuntimeAdapter(
    private val webView: WebView,
    private val core: Core,
) : CoroutineScope {

    override val coroutineContext = SupervisorJob()

    @JavascriptInterface
    fun log(var1: String) = core.log(var1)

    @JavascriptInterface
    fun sleep(var1: Long) = promise { core.sleep(var1) }

    @JavascriptInterface
    fun connAccept() = promise { core.connAccept() }

    @JavascriptInterface
    fun connClose(var1: String) = promise { core.connClose(var1) }

    @JavascriptInterface
    fun connRead(var1: String, var2: Long) = promise { core.connRead(var1, var2) }

    @JavascriptInterface
    fun connWrite(var1: String, var2: ByteArray) = promise { core.connWrite(var1, var2) }

    @JavascriptInterface
    fun connReadLn(var1: String) = promise { core.connReadLn(var1) }

    @JavascriptInterface
    fun connWriteLn(var1: String, var2: String) = promise { core.connWriteLn(var1, var2) }

    @JavascriptInterface
    fun query(var1: String?, var2: String) = promise { core.queryString(var1, var2) }

    @JavascriptInterface
    fun resolve(var1: String) = promise { core.resolve(var1) }

    @JavascriptInterface
    fun serviceClose() = promise { core.serviceClose() }

    @JavascriptInterface
    fun serviceRegister() = promise { core.serviceRegister() }

    @JavascriptInterface
    fun nodeInfo(var1: String) = promise { core.nodeInfoString(var1) }

    private fun <T> promise(block: suspend CoroutineScope.(UUID) -> T): String {
        val id = UUID.randomUUID()
        launch(Dispatchers.IO) {
            val fn = try {
                val result = block(id).takeIf { it != Unit }
                val quoted = JSONObject.quote(result?.toString())
                "window._resolve(\"$id\", $quoted)"
            } catch (e: Throwable) {
                val message = e.message
                "window._reject(\"$id\", \"$message\")"
            }
            launch(Dispatchers.Main) {
                webView.evaluateJavascript(fn) {
                    Log.d(Tag, "done: $fn, $it")
                }
            }
        }
        return id.toString()
    }

    companion object {
        private val Tag = HtmlRuntimeAdapter::class.java.simpleName
    }
}
