package cc.cryptopunks.portal.html

import android.annotation.SuppressLint
import android.content.ComponentName
import android.content.Context
import android.content.Intent
import android.os.Bundle
import android.webkit.WebChromeClient
import android.webkit.WebView
import androidx.activity.ComponentActivity
import androidx.activity.addCallback
import cc.cryptopunks.portal.core.mobile.App
import cc.cryptopunks.portal.core.mobile.Runtime
import org.koin.android.ext.android.inject

internal fun Context.startHtmlAppActivity(src: String, slot: Int) {
    val intent = Intent()
    val activity = "$packageName.js.JsAppActivity\$Id$slot"
    intent.component = ComponentName(packageName, activity)
    intent.putExtra("src", src)
    startActivity(intent)
}

sealed class HtmlAppActivity : ComponentActivity() {

    private val webView: WebView by lazy { WebView(this) }
    private val runtime: Runtime by inject()
    private val registry: HtmlAppsRegistry by inject()
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

    class Id0 : HtmlAppActivity()
    class Id1 : HtmlAppActivity()
    class Id2 : HtmlAppActivity()
    class Id3 : HtmlAppActivity()
    class Id4 : HtmlAppActivity()
    class Id5 : HtmlAppActivity()
    class Id6 : HtmlAppActivity()
    class Id7 : HtmlAppActivity()
    class Id8 : HtmlAppActivity()
    class Id9 : HtmlAppActivity()
    class Id10 : HtmlAppActivity()
    class Id11 : HtmlAppActivity()
    class Id12 : HtmlAppActivity()
    class Id13 : HtmlAppActivity()
    class Id14 : HtmlAppActivity()
    class Id15 : HtmlAppActivity()
    class Id16 : HtmlAppActivity()
    class Id17 : HtmlAppActivity()
    class Id18 : HtmlAppActivity()
    class Id19 : HtmlAppActivity()
}

