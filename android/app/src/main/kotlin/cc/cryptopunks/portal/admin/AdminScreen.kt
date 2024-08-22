package cc.cryptopunks.portal.admin

import androidx.compose.foundation.gestures.scrollBy
import androidx.compose.foundation.horizontalScroll
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.selection.SelectionContainer
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.Icon
import androidx.compose.material.IconButton
import androidx.compose.material.OutlinedTextField
import androidx.compose.material.Text
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Send
import androidx.compose.material.minimumInteractiveComponentSize
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import kotlinx.coroutines.launch
import org.koin.androidx.compose.navigation.koinNavViewModel
import org.koin.compose.koinInject

@Composable
fun AdminScreen(
    model: AdminViewModel = koinNavViewModel(),
    preferences: AdminPreferences = koinInject(),
) {
    val output by model.output.collectAsStateWithLifecycle()

    AdminScreen(
        output = output,
        query = model.query,
        softWrap = preferences.softWrap.value,
        queryChanged = { model.query = it },
        queryTrigger = { model.query() },
    )
}

@Preview
@Composable
private fun AdminScreenPreview() {
    var output by remember {
        mutableStateOf("test output\n")
    }
    var query by remember {
        mutableStateOf("")
    }
    AdminScreen(
        output = output,
        query = query,
        queryChanged = { query = it },
        queryTrigger = {
            output += query + "\n"
            query = ""
        }
    )
}

@Composable
private fun AdminScreen(
    output: String,
    query: String,
    softWrap: Boolean = true,
    queryChanged: (String) -> Unit,
    queryTrigger: () -> Unit,
) {
    val scrollState = rememberScrollState()
    val scope = rememberCoroutineScope()
    if (!scrollState.canScrollForward) LaunchedEffect(output.length) {
        scope.launch {
            scrollState.scrollBy(scrollState.scrollBy(Float.MAX_VALUE))
        }
    }
    LaunchedEffect(Unit) {
        scrollState.scrollBy(scrollState.scrollBy(Float.MAX_VALUE))
    }

    Column {
        Box(
            modifier = when {
                softWrap -> Modifier
                else -> Modifier
                    .fillMaxSize()
                    .horizontalScroll(rememberScrollState())
            }.weight(1f),
        ) {
            SelectionContainer {
                Text(
                    text = output,
                    modifier = Modifier
                        .fillMaxSize()
                        .verticalScroll(scrollState),
                    fontFamily = FontFamily.Monospace,
                    softWrap = softWrap,
                )
            }
        }
        Box {
            OutlinedTextField(
                value = query,
                onValueChange = queryChanged,
                modifier = Modifier
                    .verticalScroll(rememberScrollState(), reverseScrolling = true)
                    .fillMaxWidth(),
                trailingIcon = {
                    Spacer(
                        modifier = Modifier
                            .minimumInteractiveComponentSize()
                            .size(24.dp)
                    )
                },
                placeholder = {
                    Text(text = "type command here...")
                },
            )
            IconButton(
                onClick = queryTrigger,
                modifier = Modifier
                    .align(Alignment.BottomEnd)
                    .padding(bottom = 4.dp)
            ) {
                Icon(
                    imageVector = Icons.Default.Send,
                    contentDescription = ""
                )
            }
        }
    }
}
