package cc.cryptopunks.portal.html

import android.util.Log
import android.webkit.WebResourceError
import android.webkit.WebResourceRequest
import android.webkit.WebResourceResponse
import android.webkit.WebView
import android.webkit.WebViewClient
import cc.cryptopunks.portal.core.mobile.Assets
import cc.cryptopunks.portal.ext.mobile.inputStream
import kotlin.math.min

internal class PortalWebViewClient(
    private val assets: Assets
) : WebViewClient() {
    private val tag = javaClass.simpleName

    override fun shouldInterceptRequest(
        view: WebView,
        request: WebResourceRequest
    ): WebResourceResponse? {
        return when (request.url.scheme) {
            "file" -> {
                Log.d(tag, "requesting url: ${request.url}")

                val host = request.url.host ?: return null
                var path = request.url.path ?: return null
                path = if (path == "/") host else path.trimStart('/')

                val asset = assets.get(path)
                val mime = asset.mime()
                val encoding = asset.encoding()
                Log.d(tag, "requesting: $mime, $encoding, $path")

                return WebResourceResponse(
                    mime,
                    encoding,
                    asset.data().inputStream(),
                )
            }

            else -> super.shouldInterceptRequest(view, request)
        }
    }

    override fun onReceivedError(
        view: WebView,
        request: WebResourceRequest,
        error: WebResourceError,
    ) {
        Log.d(
            "PortalWebViewClient",
            "error: ${request.method} ${request.url}: ${error.errorCode} ${error.description}"
        )
    }
}
