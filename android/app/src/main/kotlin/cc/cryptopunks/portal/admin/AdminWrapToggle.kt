package cc.cryptopunks.portal.admin

import androidx.compose.runtime.Composable
import cc.cryptopunks.portal.compose.SoftWrapToggle
import org.koin.compose.koinInject

@Composable
fun AdminWrapToggle(
    preferences: AdminPreferences = koinInject(),
) = SoftWrapToggle(
    state = preferences.softWrap,
)
