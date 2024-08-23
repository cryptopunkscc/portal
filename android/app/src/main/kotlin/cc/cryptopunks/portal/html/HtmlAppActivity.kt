package cc.cryptopunks.portal.html

import android.annotation.SuppressLint
import android.content.Intent
import android.os.Bundle
import android.webkit.WebChromeClient
import android.webkit.WebView
import androidx.activity.ComponentActivity
import androidx.activity.addCallback
import cc.cryptopunks.portal.core.mobile.App
import cc.cryptopunks.portal.core.mobile.Runtime
import org.koin.android.ext.android.inject

sealed class HtmlAppActivity(val slot: Int) : ComponentActivity() {

    private val webView: WebView by lazy { WebView(this) }
    private val runtime: Runtime by inject()
    private val registry: HtmlAppRegistry by inject()
    private var app: App? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        registry.onCreate(this)
        setContentView(webView)
        setupWebView()
        setupOnBackPressedCallback()
    }

    private fun reloadApp(intent: Intent = getIntent()): App {
        val src = intent.getStringExtra("src")
        app?.run { runtime().close() }
        return runtime.app(src).also { app ->
            this.app = app
            title = app.manifest().title
        }
    }

    private fun setupWebView(app: App = reloadApp()) = webView.apply {
        val assets = app.assets()
        webViewClient = PortalWebViewClient(assets)
        webChromeClient = WebChromeClient()
        settings.apply {
            @SuppressLint("SetJavaScriptEnabled")
            javaScriptEnabled = true
            domStorageEnabled = true
            useWideViewPort = true
            allowFileAccess = true
            allowContentAccess = true
            domStorageEnabled = true
            databaseEnabled = true
            javaScriptCanOpenWindowsAutomatically = true
        }
        addJavascriptInterface(
            HtmlRuntimeAdapter(webView, app.runtime()),
            "_app_host"
        )
        assets.get("index.html").run {
            loadData(
                data().readAll().decodeToString(),
                mime(),
                encoding(),
            )
        }
    }

    private fun setupOnBackPressedCallback() {
        onBackPressedDispatcher.addCallback(this) {
            if (webView.canGoBack()) {
                webView.goBack()
            } else {
                isEnabled = false
                onBackPressedDispatcher.onBackPressed()
                isEnabled = true
            }
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        registry.onDestroy(this)
        webView.destroy()
        app?.run { runtime().close() }
    }

    companion object {
        val slots: List<Int> = (0..<HtmlAppActivity::class.nestedClasses.size).toList()
    }

    class Slot0 : HtmlAppActivity(0)
    class Slot1 : HtmlAppActivity(1)
    class Slot2 : HtmlAppActivity(2)
    class Slot3 : HtmlAppActivity(3)
    class Slot4 : HtmlAppActivity(4)
    class Slot5 : HtmlAppActivity(5)
    class Slot6 : HtmlAppActivity(6)
    class Slot7 : HtmlAppActivity(7)
    class Slot8 : HtmlAppActivity(8)
    class Slot9 : HtmlAppActivity(9)
    class Slot10 : HtmlAppActivity(10)
    class Slot11 : HtmlAppActivity(11)
    class Slot12 : HtmlAppActivity(12)
    class Slot13 : HtmlAppActivity(13)
    class Slot14 : HtmlAppActivity(14)
    class Slot15 : HtmlAppActivity(15)
    class Slot16 : HtmlAppActivity(16)
    class Slot17 : HtmlAppActivity(17)
    class Slot18 : HtmlAppActivity(18)
    class Slot19 : HtmlAppActivity(19)
}

