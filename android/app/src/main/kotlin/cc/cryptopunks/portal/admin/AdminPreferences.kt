package cc.cryptopunks.portal.admin

import android.content.Context
import android.content.Context.MODE_PRIVATE
import cc.cryptopunks.portal.compose.mutableStateOf

class AdminPreferences(
    context: Context,
) {
    private val sharedPreferences = context.getSharedPreferences("admin", MODE_PRIVATE)

    var softWrap = sharedPreferences.mutableStateOf("soft_wrap", true)
}
