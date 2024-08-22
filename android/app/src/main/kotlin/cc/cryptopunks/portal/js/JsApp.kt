package cc.cryptopunks.portal.js

import android.content.Context

internal data class JsApp(
    val name: String,
    val dir: String,
    val description: String = "",
    val icon: String? = null,
    val service: String? = null,
    val activity: Int = -1,
)

internal val Context.appsDir get() = dataDir.resolve("apps")
