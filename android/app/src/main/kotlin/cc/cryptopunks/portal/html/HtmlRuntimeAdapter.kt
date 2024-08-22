package cc.cryptopunks.portal.html

import android.util.Log
import android.webkit.JavascriptInterface
import android.webkit.WebView
import cc.cryptopunks.portal.core.bind.Runtime
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.launch
import org.json.JSONObject
import java.util.UUID

internal class HtmlRuntimeAdapter(
    private val webView: WebView,
    private val runtime: Runtime,
) : CoroutineScope {

    override val coroutineContext = SupervisorJob()

    @JavascriptInterface
    fun log(var1: String) = runtime.log(var1)

    @JavascriptInterface
    fun sleep(var1: Long) = promise { runtime.sleep(var1) }

    @JavascriptInterface
    fun connAccept(var1: String) = promise { runtime.connAccept(var1) }

    @JavascriptInterface
    fun connClose(var1: String) = promise { runtime.connClose(var1) }

    @JavascriptInterface
    fun connRead(var1: String) = promise { runtime.connRead(var1) }

    @JavascriptInterface
    fun connWrite(var1: String, var2: String) = promise { runtime.connWrite(var1, var2) }

    @JavascriptInterface
    fun query(var1: String?, var2: String) = promise { runtime.query(var1, var2) }

    @JavascriptInterface
    fun queryName(var1: String, var2: String) = promise { runtime.queryName(var1, var2) }

    @JavascriptInterface
    fun resolve(var1: String) = promise { runtime.resolve(var1) }

    @JavascriptInterface
    fun serviceClose(var1: String) = promise { runtime.serviceClose(var1) }

    @JavascriptInterface
    fun serviceRegister(var1: String) = promise { runtime.serviceRegister(var1) }

    @JavascriptInterface
    fun nodeInfo(var1: String) = promise {
        runtime.nodeInfo(var1)?.run {
            """{"name":"$name","identity":"$identity"}"""
        } ?: "{}"
    }

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
