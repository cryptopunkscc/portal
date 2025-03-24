package cc.cryptopunks.portal.exception

import androidx.compose.foundation.horizontalScroll
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.FloatingActionButton
import androidx.compose.material.Icon
import androidx.compose.material.Scaffold
import androidx.compose.material.Text
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Close
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.compose.ui.window.Dialog
import androidx.compose.ui.window.DialogProperties
import cc.cryptopunks.portal.compose.AstralTheme

@Preview
@Composable
private fun StackTracePreview() = AstralTheme {
    val message = "Test error preview, with some long description, that will overflow in dialog, so the UI can be adjusted."
    StackTraceView(err = Throwable(message))
}

@Composable
internal fun StackTraceView(
    err: Throwable,
    dismiss: () -> Unit = {}
) {
    Dialog(
        onDismissRequest = dismiss,
        properties = DialogProperties(usePlatformDefaultWidth = false)
    ) {
        Scaffold(
            floatingActionButton = {
                FloatingActionButton(onClick = dismiss) {
                    Icon(imageVector = Icons.Default.Close, contentDescription = "cancel")
                }
            }
        ) { paddingValues ->
            Box(
                modifier = Modifier
                    .padding(paddingValues)
                    .verticalScroll(rememberScrollState())
            ) {
                Text(
                    modifier = Modifier
                        .padding(bottom = 96.dp)
                        .horizontalScroll(rememberScrollState()),
                    text = err.stackTraceToString(),
                    fontSize = 10.sp,
                )
            }
        }
    }
}
