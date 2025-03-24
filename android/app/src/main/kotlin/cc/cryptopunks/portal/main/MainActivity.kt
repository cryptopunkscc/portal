package cc.cryptopunks.portal.main

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Text
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.text.TextStyle
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.lifecycleScope
import cc.cryptopunks.portal.LAUNCHER
import cc.cryptopunks.portal.StartHtmlApp
import cc.cryptopunks.portal.compose.inject
import cc.cryptopunks.portal.core.mobile.Event
import cc.cryptopunks.portal.core.mobile.Mobile
import cc.cryptopunks.portal.exception.ExceptionsState
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.flow.filter
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.mapNotNull
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
        setContent {
            inject.Errors()
            Column(
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.Center,
            ) {
                val event by events.collectAsStateWithLifecycle()
                Text(
                    text = when (event.msg) {
                        Mobile.STARTING -> "Portal starting..."
                        Mobile.STARTED -> "Portal started."
                        Mobile.STOPPED -> "Portal stopped."
                        else -> "Invalid status"
                    },
                    style = MaterialTheme.typography.h2
                )
            }
        }
    }

    override fun onResume() {
        super.onResume()
        if (!mainPermissions.ask(this)) {
            startAstralService()
        }
    }
}
