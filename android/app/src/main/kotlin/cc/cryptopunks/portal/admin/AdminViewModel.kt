package cc.cryptopunks.portal.admin

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import androidx.lifecycle.ViewModel
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch

class AdminViewModel(
    private val client: AdminClient,
    private val scope: CoroutineScope,
) : ViewModel() {

    val output = client.output

    var query by mutableStateOf("")

    fun query() = scope.launch {
        client.query(query)
    }
}
