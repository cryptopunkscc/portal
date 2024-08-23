package cc.cryptopunks.portal.main

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.material.Text
import androidx.lifecycle.lifecycleScope
import cc.cryptopunks.portal.LAUNCHER
import cc.cryptopunks.portal.StartHtmlApp
import cc.cryptopunks.portal.compose.inject
import cc.cryptopunks.portal.core.mobile.Mobile
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import org.koin.android.ext.android.inject
import org.koin.androidx.viewmodel.ext.android.viewModel

class MainActivity : ComponentActivity() {
    private val mainPermissions: MainPermissions by viewModel()
    private val startHtmlApp: StartHtmlApp by inject()
    private val events: MainEvents by inject()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        lifecycleScope.launch {
            events.first { it.msg == Mobile.STARTED }
            startHtmlApp(LAUNCHER)
            finish()
        }
        if (events.value.msg != Mobile.STARTED) setContent {
            Text(text = "Portal starting...")
            inject.Errors()
        }
    }

    override fun onResume() {
        super.onResume()
        if (!mainPermissions.ask(this)) {
            startAstralService()
        }
    }
}
