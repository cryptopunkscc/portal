package cc.cryptopunks.portal

import android.content.Context
import android.content.Intent
import android.content.Intent.ACTION_VIEW
import android.content.pm.PackageManager
import android.net.Uri
import androidx.core.content.ContextCompat

object Permissions {

    object Key {
        const val Message = "message"
        const val Required = "request"
        const val Rejected = "rejected"
    }

    fun request(
        message: String,
        vararg required: String,
    ) = Intent(ACTION_VIEW, Uri.parse("astral://permissions")).apply {
        putExtra(Key.Message, message)
        putExtra(Key.Required, required)
    }

    fun result(
        rejected: Array<String>,
    ) = Intent().apply {
        putExtra(Key.Rejected, rejected)
    }

    fun getMessage(intent: Intent): String =
        intent.getStringExtra(Key.Message) ?: ""

    fun getRequired(intent: Intent): Array<String> =
        intent.getStringArrayExtra(Key.Required) ?: emptyArray()

    fun getRejected(intent: Intent): Array<String> =
        intent.getStringArrayExtra(Key.Rejected) ?: emptyArray()
}

fun Context.hasPermissions(name: String): Boolean =
    ContextCompat.checkSelfPermission(applicationContext, name) ==
        PackageManager.PERMISSION_GRANTED
