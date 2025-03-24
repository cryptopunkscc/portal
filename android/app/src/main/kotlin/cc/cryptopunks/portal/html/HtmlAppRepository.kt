package cc.cryptopunks.portal.html

import android.content.ComponentName
import android.content.Context
import android.content.Intent
import cc.cryptopunks.portal.Activities
import cc.cryptopunks.portal.StartHtmlApp
import java.lang.ref.WeakReference

class HtmlAppRepository(
    private val activities: Activities,
) : HtmlAppRegistry, StartHtmlApp {

    private val appActivities = mutableListOf<HtmlAppActivity>()

    private val usedSlots get() = appActivities.map(HtmlAppActivity::slot).toSet()
    private val freeSlots get() = HtmlAppActivity.slots - usedSlots
    private val nextSlot get() = freeSlots.first()

    override var current: HtmlAppActivity? = appActivities.lastOrNull()

    override fun onCreate(activity: HtmlAppActivity) {
        appActivities.add(activity)
    }

    override fun onDestroy(activity: HtmlAppActivity) {
        appActivities.remove(activity)
    }

    override fun invoke(src: String) {
        activities.last().startHtmlAppActivity(src, nextSlot)
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
    "$packageName.html.HtmlAppActivity\$Slot$slot"
