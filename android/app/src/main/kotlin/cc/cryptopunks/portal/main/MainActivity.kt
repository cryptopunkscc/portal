package cc.cryptopunks.portal.main

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import cc.cryptopunks.portal.compose.inject
import org.koin.androidx.viewmodel.ext.android.viewModel

class MainActivity : ComponentActivity() {

    private val mainPermissions: MainPermissions by viewModel()
    private val startAstralService by lazy { startAstralService() }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MainScreen()
            inject.Errors()
        }
    }

    override fun onResume() {
        super.onResume()
        if (!mainPermissions.ask(this))
            startAstralService
    }
}
