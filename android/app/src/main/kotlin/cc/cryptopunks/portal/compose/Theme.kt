package cc.cryptopunks.portal.compose

import androidx.compose.material.MaterialTheme
import androidx.compose.material.darkColors
import androidx.compose.runtime.Composable

@Composable
fun AstralTheme(
    content: @Composable () -> Unit,
) = MaterialTheme(
    colors = darkColors(),
    content = content
)
