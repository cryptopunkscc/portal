package cc.cryptopunks.portal.main

import android.annotation.SuppressLint
import android.content.Context
import android.os.Build
import android.os.PowerManager
import android.provider.Settings
import androidx.core.content.getSystemService
import androidx.lifecycle.ViewModel
import cc.cryptopunks.portal.Permissions
import cc.cryptopunks.portal.hasPermissions

internal class MainPermissions : ViewModel() {

    @SuppressLint("BatteryLife")
    private val remaining = buildList {
        add(
            Data(
                Settings.ACTION_REQUEST_IGNORE_BATTERY_OPTIMIZATIONS,
                "Turning off battery optimization, to keep connections alive in background",
            ) {
                getSystemService<PowerManager>()?.isIgnoringBatteryOptimizations(packageName) == true
            }
        )
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) add(
            Data(
                android.Manifest.permission.POST_NOTIFICATIONS,
                "Allow service notification, to keep connections alive in background",
            )
        )
    }.iterator()

    fun ask(context: Context) : Boolean {
        while (remaining.hasNext()) {
            val (permission, message, granted) = remaining.next()
            if (!granted(context)) {
                context.startActivity(Permissions.request(message, permission))
                return true
            }
        }
        return false
    }

    private data class Data(
        val permission: String,
        val message: String,
        val check: Context.() -> Boolean = { hasPermissions(permission) },
    )
}
