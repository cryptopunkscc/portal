package cc.cryptopunks.portal.js

import android.util.Log
import android.webkit.MimeTypeMap
import android.webkit.WebResourceRequest
import android.webkit.WebResourceResponse
import android.webkit.WebView
import android.webkit.WebViewClient
import cc.cryptopunks.portal.core.mobile.Assets
import cc.cryptopunks.portal.core.mobile.Reader
import cc.cryptopunks.portal.ext.mobile.inputStream
import java.io.File
import java.io.InputStream

internal class FileResolvingWebViewClient(
    private val dir: File,
) : WebViewClient() {

    private val tag = javaClass.simpleName

    override fun shouldInterceptRequest(
        view: WebView,
        request: WebResourceRequest,
    ): WebResourceResponse? {
        when (request.url.scheme) {
            "file" -> {

                // Resolve requested file
                val path = request.url.path ?: return null
                Log.d(tag, "requesting path: $path")
                val relative = if (path.startsWith('/')) path.drop(1) else path
                val file = dir.resolve(relative)
                Log.d(tag, "resolved file: ${file.absolutePath}, exist: ${file.exists()}")
                file.exists() || return null

                // Prepare response data
                val mimeType = MimeTypeMap.getSingleton().getMimeTypeFromExtension(file.extension)
                val encoding = "UTF-8"
                val data = file.inputStream()

                return WebResourceResponse(mimeType, encoding, data)
            }

            else -> {
                return super.shouldInterceptRequest(view, request)
            }
        }
    }
}
