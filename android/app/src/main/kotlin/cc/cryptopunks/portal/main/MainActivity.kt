package cc.cryptopunks.portal.main

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.lifecycleScope
import cc.cryptopunks.portal.LAUNCHER
import cc.cryptopunks.portal.StartHtmlApp
import cc.cryptopunks.portal.Status
import cc.cryptopunks.portal.compose.AstralTheme
import cc.cryptopunks.portal.compose.inject
import cc.cryptopunks.portal.core.mobile.Mobile
import cc.cryptopunks.portal.onboarding.OnBoardingScreen
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import org.koin.android.ext.android.inject
import org.koin.androidx.viewmodel.ext.android.viewModel
import org.koin.compose.koinInject

class MainActivity : ComponentActivity() {
    private val mainPermissions: MainPermissions by viewModel()
    private val startHtmlApp: StartHtmlApp by inject()
    private val status: Status by inject()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        lifecycleScope.launch {
            status.first { it == Mobile.STARTED }
            startHtmlApp(LAUNCHER)
            finishAndRemoveTask()
        }
        setContent {
            MainScreen()
        }
    }

    override fun onResume() {
        super.onResume()
        if (!mainPermissions.ask(this)) {
            startAstralService()
        }
    }
}

@Composable
fun MainScreen() = AstralTheme {
    inject.Errors()
    val status by koinInject<Status>().collectAsStateWithLifecycle()
    when (status) {
        Mobile.FRESH -> OnBoardingScreen()
//        Mobile.STARTED -> HtmlAppScreen(src = LAUNCHER)
        else -> Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center,
        ) {
            Text(
                text = when (status) {
                    Mobile.STARTING -> "Portal starting..."
                    Mobile.STARTED -> "Portal started."
                    Mobile.STOPPED -> "Portal stopped."
//                    Mobile.FRESH -> "Not configured."
                    else -> "Unknown status"
                },
                style = MaterialTheme.typography.h2
            )
        }
    }
}