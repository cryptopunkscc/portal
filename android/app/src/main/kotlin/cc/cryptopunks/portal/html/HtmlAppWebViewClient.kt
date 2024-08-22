package cc.cryptopunks.portal.html

import android.util.Log
import android.webkit.WebResourceRequest
import android.webkit.WebResourceResponse
import android.webkit.WebView
import android.webkit.WebViewClient
import cc.cryptopunks.portal.core.mobile.Assets
import cc.cryptopunks.portal.ext.mobile.inputStream


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
                val path = request.url.path ?: return null
                Log.d(tag, "requesting path: $path")
                val asset = assets.get(path)

                return WebResourceResponse(
                    asset.mime(),
                    asset.encoding(),
                    asset.data().inputStream(),
                )
            }

            else -> super.shouldInterceptRequest(view, request)
        }
    }
}
