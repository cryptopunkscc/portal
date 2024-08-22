package cc.cryptopunks.portal.js

import android.annotation.SuppressLint
import android.content.ComponentName
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.webkit.WebChromeClient
import android.webkit.WebView
import androidx.activity.ComponentActivity
import androidx.activity.addCallback
import cc.cryptopunks.portal.core.mobile.App
import cc.cryptopunks.portal.core.mobile.Runtime
import org.koin.android.ext.android.inject

internal fun Context.startJsAppActivity(app: JsApp) {
    val dir = appsDir.resolve(app.dir)
    val file = dir.resolve("index.html")
    val uri = Uri.fromFile(file)
    val intent = Intent()
    val activity = "$packageName.js.JsAppActivity\$Id${app.activity}"
    intent.component = ComponentName(packageName, activity)
    intent.data = uri
    intent.putExtra("title", app.name)
    startActivity(intent)
}

sealed class JsAppActivity : ComponentActivity() {

    private val webView: WebView by lazy { WebView(this) }
    private val runtime: Runtime by inject()
    private var app: App? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
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
            WebViewRuntimeAdapter(webView, app.runtime()),
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
        webView.destroy()
        app?.run { runtime().close() }
    }

    class Id0 : JsAppActivity()
    class Id1 : JsAppActivity()
    class Id2 : JsAppActivity()
    class Id3 : JsAppActivity()
    class Id4 : JsAppActivity()
    class Id5 : JsAppActivity()
    class Id6 : JsAppActivity()
    class Id7 : JsAppActivity()
    class Id8 : JsAppActivity()
    class Id9 : JsAppActivity()
    class Id10 : JsAppActivity()
    class Id11 : JsAppActivity()
    class Id12 : JsAppActivity()
    class Id13 : JsAppActivity()
    class Id14 : JsAppActivity()
    class Id15 : JsAppActivity()
    class Id16 : JsAppActivity()
    class Id17 : JsAppActivity()
    class Id18 : JsAppActivity()
    class Id19 : JsAppActivity()

    companion object {
        const val Limit = 20
    }
}

