package cc.cryptopunks.portal

import android.app.Activity
import cc.cryptopunks.portal.core.mobile.Event
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.StateFlow

const val LAUNCHER = "launcher"

fun interface StartHtmlApp : (String) -> Unit

interface CoreEvents : StateFlow<Event>

interface Activities : List<Activity>
