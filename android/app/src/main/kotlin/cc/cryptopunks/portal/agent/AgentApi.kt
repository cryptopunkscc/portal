package cc.cryptopunks.portal.agent

import android.content.Context
import cc.cryptopunks.portal.StartHtmlApp
import cc.cryptopunks.portal.core.mobile.Api
import cc.cryptopunks.portal.core.mobile.Event
import cc.cryptopunks.portal.core.mobile.Net
import cc.cryptopunks.portal.main.MainEvents

class AgentApi(
    private val context: Context,
    private val events: MainEvents,
    private val startHtmlApp: StartHtmlApp,
) : Api {
    override fun cacheDir(): String = context.cacheDir.path

    override fun dataDir(): String = context.dataDir.path

    override fun dbDir(): String = context.getDatabasePath("node").parent.orEmpty()

    override fun startHtml(src: String, arg: String) = startHtmlApp(src)

    override fun event(event: Event) {
        events.tryEmit(event)
    }

    override fun net(): Net = AgentNet()
}

