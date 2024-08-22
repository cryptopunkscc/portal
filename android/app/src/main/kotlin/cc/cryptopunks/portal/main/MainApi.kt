package cc.cryptopunks.portal.main

import android.content.Context
import cc.cryptopunks.portal.core.mobile.Api
import cc.cryptopunks.portal.core.mobile.Event
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.asSharedFlow
import java.io.File

class MainApi(
    private val context: Context,
    private val getCurrentActivity: GetCurrentActivity,
) : Api {

    private val _events = MutableSharedFlow<Event>(replay = 1, extraBufferCapacity = 32)
    val events = _events.asSharedFlow()

    override fun event(event: Event) {
        _events.tryEmit(event)
    }

    override fun nodeRoot(): String = context.dataDir
        .resolve("astrald")
        .apply(File::mkdirs)
        .path

    override fun requestHtml(src: String) {
        val activity = requireNotNull(getCurrentActivity())

        activity.startActivity()
    }
}