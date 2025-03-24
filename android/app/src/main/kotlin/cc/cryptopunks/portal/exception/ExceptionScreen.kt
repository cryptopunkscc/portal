package cc.cryptopunks.portal.exception

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.Button
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Surface
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalClipboardManager
import androidx.compose.ui.text.AnnotatedString
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.window.Dialog
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import cc.cryptopunks.portal.compose.AstralTheme
import org.koin.compose.koinInject

@Preview
@Composable
private fun ErrorsPreview() = AstralTheme {
    val message = "Test error preview, with some long description, that will overflow in dialog, so the UI can be adjusted."
    val errors = ExceptionsState()
    errors.plusAssign(Exception(message))
    ErrorsScreen(errors = errors)
}

@Composable
fun ErrorsScreen(
    errors: ExceptionsState = koinInject(),
) {
    val error by errors.current.collectAsStateWithLifecycle(null)
    error?.let { err ->
        val drop: () -> Unit = {
            errors.pop()
        }
        var stacktrace by remember(err) {
            mutableStateOf(false)
        }
        val clipboard = LocalClipboardManager.current
        Dialog(drop) {
            Surface(
                shape = RoundedCornerShape(16.dp),
                elevation = 4.dp,
            ) {
                Column(Modifier.padding(16.dp)) {
                    Text(
                        text = "Unexpected Error",
                        modifier = Modifier
                            .align(Alignment.CenterHorizontally)
                            .padding(16.dp),
                        style = MaterialTheme.typography.h5,
                    )
                    err.message?.let { message ->
                        Text(
                            text = message,
                            modifier = Modifier
                                .align(Alignment.CenterHorizontally)
                                .padding(16.dp),
                            style = MaterialTheme.typography.subtitle1,
                            textAlign = TextAlign.Center
                        )
                        Spacer(modifier = Modifier.height(24.dp))
                    }
                    Row(
                        horizontalArrangement = Arrangement.SpaceAround,
                        modifier = Modifier.fillMaxWidth(),
                    ) {
                        Button(onClick = {
                            stacktrace = !stacktrace
                        }) {
                            Text(text = "trace")
                        }
                        Button(onClick = {
                            val string = err.stackTraceToString()
                            clipboard.setText(AnnotatedString(string))
                        }) {
                            Text(text = "copy")
                        }
                        Button(onClick = drop) {
                            Text(text = "close")
                        }
                    }
                }
            }
        }
        if (stacktrace) StackTraceView(err = err) {
            stacktrace = false
        }
    }
}
