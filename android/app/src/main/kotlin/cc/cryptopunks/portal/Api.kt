package cc.cryptopunks.portal

import android.app.Activity
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.StateFlow

const val LAUNCHER = "portal.launcher"

fun interface StartHtmlApp : (String) -> Unit

interface Status : StateFlow<Int>

interface Errors : Flow<String>

interface Activities : List<Activity>
