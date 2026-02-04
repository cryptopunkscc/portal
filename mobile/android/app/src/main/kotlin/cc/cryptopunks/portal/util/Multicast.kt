package cc.cryptopunks.portal.util

import android.content.Context
import android.net.wifi.WifiManager
import androidx.core.content.getSystemService

internal fun Context.acquireMulticastWakeLock(tag: String): WifiManager.MulticastLock =
    applicationContext.getSystemService<WifiManager>()!!
        .createMulticastLock(tag).apply {
            setReferenceCounted(true)
        }
