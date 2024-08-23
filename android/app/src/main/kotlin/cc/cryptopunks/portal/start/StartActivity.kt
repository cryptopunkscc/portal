package cc.cryptopunks.portal.start

import android.app.Activity
import android.os.Bundle
import cc.cryptopunks.portal.StartHtmlApp
import org.koin.android.ext.android.inject

class StartActivity : Activity() {

    private val startHtmlApp: StartHtmlApp by inject()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        startHtmlApp("launcher")
    }
}