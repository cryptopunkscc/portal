package cc.cryptopunks.portal.html

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.addCallback
import org.koin.android.ext.android.inject

sealed class HtmlAppActivity(val slot: Int) : ComponentActivity() {

    private val webView: HtmlAppView by lazy { HtmlAppView(this) }
    private val repository: HtmlAppRepository by inject()
    val src by lazy { intent.getStringExtra("src")!! }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        prepare() || return
        setContentView(webView)
        webView.loadApp(src)
        setupOnBackPressedCallback()
    }

    private fun prepare(): Boolean = when {
        runCatching { src }.isFailure -> false
        this in repository -> false.also { repository(src) }
        else -> true
    }

    private fun setupOnBackPressedCallback() {
        onBackPressedDispatcher.addCallback(this) {
            if (webView.canGoBack()) webView.goBack()
            else finishAndRemoveTask()
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        webView.destroy()
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

