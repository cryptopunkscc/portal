package cc.cryptopunks.portal.html

import android.annotation.SuppressLint
import android.content.Context
import android.webkit.WebChromeClient
import android.webkit.WebView
import androidx.compose.runtime.Composable
import androidx.compose.ui.viewinterop.AndroidView
import cc.cryptopunks.portal.core.mobile.App
import cc.cryptopunks.portal.core.mobile.Core
import org.koin.core.component.KoinComponent
import org.koin.core.component.inject

class HtmlAppView(context: Context) : WebView(context), KoinComponent {

    private val portalWebViewClient = PortalWebViewClient()
    private val core: Core by inject()
    private var app: App? = null

    init {
        webViewClient = portalWebViewClient
        webChromeClient = WebChromeClient()
        settings.apply {
            @SuppressLint("SetJavaScriptEnabled")
            javaScriptEnabled = true
            domStorageEnabled = true
            useWideViewPort = true
            allowFileAccess = true
//            allowFileAccessFromFileURLs = true
//            allowUniversalAccessFromFileURLs = true
            allowContentAccess = true
            domStorageEnabled = true
            databaseEnabled = true
            javaScriptCanOpenWindowsAutomatically = true
        }
    }

    fun loadApp(src: String): App {
        app?.core()?.close()

        val app = core.app(src)
        val assets = app.assets()
        val adapter = HtmlRuntimeAdapter(this, app.core())

        this.app = app
        portalWebViewClient.assets = assets
        removeJavascriptInterface("_app_host")
        addJavascriptInterface(adapter, "_app_host")
        loadUrl("file://index.html")
        return app
    }

    override fun destroy() {
        super.destroy()
        app?.core()?.close()
        app = null
        portalWebViewClient.assets = null
    }
}

@Composable
fun HtmlAppScreen(src: String) {
    AndroidView(factory = ::HtmlAppView) { view ->
        view.loadApp(src)
    }
}