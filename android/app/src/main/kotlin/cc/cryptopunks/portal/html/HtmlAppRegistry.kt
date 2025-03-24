package cc.cryptopunks.portal.html

interface HtmlAppRegistry {
    val current: HtmlAppActivity?
    fun onCreate(activity: HtmlAppActivity)
    fun onDestroy(activity: HtmlAppActivity)
}