package cc.cryptopunks.portal.compose

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.AlertDialog
import androidx.compose.material.Button
import androidx.compose.material.FloatingActionButton
import androidx.compose.material.Icon
import androidx.compose.material.IconButton
import androidx.compose.material.LocalContentAlpha
import androidx.compose.material.LocalContentColor
import androidx.compose.material.Scaffold
import androidx.compose.material.Text
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Delete
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import kotlin.random.Random

@Preview
@Composable
fun EditableItemListPreview() {
    val items = remember {
        mutableStateListOf<String>().apply {
            addAll(listOf("one", "two", "three"))
        }
    }

    EditableItemList(
        items = items,
        onAddClick = {
            items.add(Random.nextInt().toString())
        },
        onRemoveClick = {
            items.remove(it)
        }
    ) {
        Text(text = it)
    }
}

@Composable
fun <T> EditableItemList(
    items: List<T>,
    addContentDescription: String = "",
    removeContentDescription: String = "",
    onAddClick: () -> Unit = {},
    onSelectClick: (T) -> Unit = {},
    onRemoveClick: (T) -> Unit = {},
    content: @Composable (T) -> Unit,
) {
    var removeItem by remember { mutableStateOf(null as T?) }

    Scaffold(
        floatingActionButton = {
            FloatingActionButton(
                onClick = onAddClick
            ) {
                Icon(
                    imageVector = Icons.Default.Add,
                    contentDescription = addContentDescription
                )
            }
        },
    ) {
        LazyColumn(
            modifier = Modifier.padding(it)
        ) {
            items(items) { item ->
                Column {
                    Row(
                        modifier = Modifier
                            .clickable {
                                onSelectClick(item)
                            }
                            .padding(vertical = 16.dp)
                            .padding(start = 24.dp, end = 4.dp),
                        verticalAlignment = Alignment.CenterVertically,
                    ) {
                        content(item)
                        Spacer(modifier = Modifier.weight(1f))
                        IconButton(
                            onClick = {
                                removeItem = item
                            }
                        ) {
                            Icon(
                                imageVector = Icons.Default.Delete,
                                contentDescription = removeContentDescription,
                            )
                        }
                    }
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(0.5.dp)
                            .background(LocalContentColor.current.copy(alpha = LocalContentAlpha.current))
                    )
                }
            }
        }
    }

    removeItem?.let {
        val dismiss = {
            removeItem = null
        }
        AlertDialog(
            title = { Text(text = "Confirm deleting") },
            text = { content(it) },
            onDismissRequest = dismiss,
            backgroundColor = Color.DarkGray,
            confirmButton = {
                Button(onClick = {
                    dismiss()
                    onRemoveClick(it)
                }) {
                    Text(text = "OK")
                }
            },
            dismissButton = {
                Button(onClick = dismiss) {
                    Text(text = "Cancel")
                }
            },
        )
    }
}

