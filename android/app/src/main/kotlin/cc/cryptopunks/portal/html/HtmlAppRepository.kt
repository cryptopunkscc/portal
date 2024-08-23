package cc.cryptopunks.portal.html

import android.content.ComponentName
import android.content.Context
import android.content.Intent
import cc.cryptopunks.portal.StartHtmlApp
import java.lang.ref.WeakReference

class HtmlAppRepository : HtmlAppRegistry, StartHtmlApp {

    private val activities = mutableListOf<HtmlAppActivity>()
    private var currentWeakRef = WeakReference<HtmlAppActivity>(null)

    private val usedSlots get() = activities.map(HtmlAppActivity::slot).toSet()
    private val freeSlots get() = HtmlAppActivity.slots - usedSlots
    private val nextSlot get() = freeSlots.first()

    override var current: HtmlAppActivity?
        get() = currentWeakRef.get()
        set(value) {
            currentWeakRef = WeakReference(value)
        }

    override fun onCreate(activity: HtmlAppActivity) {
        activities.add(activity)
    }

    override fun onDestroy(activity: HtmlAppActivity) {
        activities.remove(activity)
    }

    override fun invoke(src: String) {
        val context = requireNotNull(current)
        context.startHtmlAppActivity(src, nextSlot)
    }
}

private fun Context.startHtmlAppActivity(src: String, slot: Int) =
    startActivity(htmlAppActivityIntent(src, slot))

private fun Context.htmlAppActivityIntent(src: String, slot: Int): Intent =
    Intent().apply {
        component = ComponentName(
            packageName,
            htmlAppActivityClassName(slot)
        )
        putExtra("src", src)
    }

private fun Context.htmlAppActivityClassName(slot: Int): String =
    "$packageName.html.JsAppActivity\$Id$slot"
