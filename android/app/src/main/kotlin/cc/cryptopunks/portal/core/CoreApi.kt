package cc.cryptopunks.portal.core

import android.content.Context
import cc.cryptopunks.portal.Errors
import cc.cryptopunks.portal.StartHtmlApp
import cc.cryptopunks.portal.Status
import cc.cryptopunks.portal.core.mobile.Api
import cc.cryptopunks.portal.core.mobile.Mobile
import cc.cryptopunks.portal.core.mobile.Net
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow

class CoreApi(
    private val context: Context,
    private val status: CoreStatus,
    private val errors: CoreErrors,
    private val startHtmlApp: StartHtmlApp,
) : Api {
    override fun cacheDir(): String = context.cacheDir.path

    override fun dataDir(): String = context.dataDir.path

    override fun dbDir(): String = context.getDatabasePath("node").parent.orEmpty()

    override fun startHtml(src: String, arg: String) = startHtmlApp(src)

    override fun status(id: Int) {
        status.tryEmit(id)
    }

    override fun error(err: String) {
        errors.tryEmit(err)
    }

    override fun net(): Net = CoreNet()
}

class CoreStatus : Status, MutableStateFlow<Int> by MutableStateFlow(Mobile.STOPPED)

class CoreErrors : Errors, MutableSharedFlow<String> by MutableSharedFlow()
