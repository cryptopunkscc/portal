package cc.cryptopunks.portal.html

interface HtmlAppsRegistry {
    fun onCreate(activity: HtmlAppActivity)
    fun onDestroy(activity: HtmlAppActivity)
}

class HtmlAppActivities : HtmlAppsRegistry {

    private val activities = mutableListOf<HtmlAppActivity>()

    override fun onCreate(activity: HtmlAppActivity) {
        activities.add(activity)
    }

    override fun onDestroy(activity: HtmlAppActivity) {
        activities.remove(activity)
    }
}