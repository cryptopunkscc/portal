package cc.cryptopunks.portal.main

import android.content.Context
import cc.cryptopunks.portal.StartHtmlApp
import cc.cryptopunks.portal.core.mobile.Api
import cc.cryptopunks.portal.core.mobile.Event
import java.io.File

class MainApi(
    private val context: Context,
    private val events: MainEvents,
    private val launchHtmlApp: StartHtmlApp,
) : Api {

    override fun event(event: Event) {
        events.tryEmit(event)
    }

    override fun nodeRoot(): String = context.dataDir
        .resolve("astrald")
        .apply(File::mkdirs)
        .path

    override fun requestHtml(src: String) = launchHtmlApp(src)
}
