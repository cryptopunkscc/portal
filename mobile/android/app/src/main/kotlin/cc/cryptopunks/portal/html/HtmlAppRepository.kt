package cc.cryptopunks.portal.html

import android.content.ComponentName
import android.content.Context
import android.content.Intent
import cc.cryptopunks.portal.Activities
import cc.cryptopunks.portal.StartHtmlApp

class HtmlAppRepository(
    private val activities: Activities,
) : StartHtmlApp {

    private val appActivities get() = activities.filterIsInstance<HtmlAppActivity>()
    private val usedSlots get() = appActivities.map(HtmlAppActivity::slot).toSet()
    private val freeSlots get() = HtmlAppActivity.slots - usedSlots
    private val nextSlot get() = freeSlots.first()

    override fun invoke(src: String) {
        activities.lastOrNull()?.startHtmlAppActivity(src, slot(src))
    }

    operator fun contains(activity: HtmlAppActivity) = appActivities.any { activity.src == it.src && it != activity }

    private fun slot(src: String) = appActivities.firstOrNull { src == it.src }?.slot ?: nextSlot
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
