package cc.cryptopunks.portal.compose

import androidx.compose.foundation.Image
import androidx.compose.material.IconButton
import androidx.compose.material.LocalContentColor
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Menu
import androidx.compose.runtime.Composable
import androidx.compose.runtime.MutableState
import androidx.compose.ui.graphics.ColorFilter
import androidx.compose.ui.tooling.preview.Preview

@Preview
@Composable
private fun SoftWrapTogglePreview() {
    SoftWrapToggle {}
}

@Composable
fun SoftWrapToggle(
    state: MutableState<Boolean>
) = SoftWrapToggle {
    state.value = !state.value
}

@Composable
fun SoftWrapToggle(
    onClick: () -> Unit,
) {
    IconButton(
        onClick = onClick,
    ) {
        Image(
            imageVector = Icons.Default.Menu,
            contentDescription = "soft wrap",
            colorFilter = ColorFilter.tint(LocalContentColor.current)
        )
    }
}
