package cc.cryptopunks.portal.html

interface HtmlAppRegistry {
    var current: HtmlAppActivity?
    fun onCreate(activity: HtmlAppActivity)
    fun onDestroy(activity: HtmlAppActivity)
}