package cc.cryptopunks.portal.main

import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.util.Log
import cc.cryptopunks.portal.core.mobile.Runtime
import org.koin.android.ext.android.inject

class MainService : Service() {

    private val tag = javaClass.simpleName
    private val runtime: Runtime by inject()

    override fun onCreate() {
        Log.d(tag, "onCreate")
        runtime.start()
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        Log.d(tag, "onStartCommand")
        startForegroundNotification()
        return START_STICKY
    }

    override fun onLowMemory() {
        Log.d(tag, "onLowMemory")
    }

    override fun onTrimMemory(level: Int) {
        Log.d(tag, "onTrimMemory")
    }

    override fun onTaskRemoved(rootIntent: Intent?) {
        Log.d(tag, "onTaskRemoved")
    }

    override fun onDestroy() {
        Log.d(tag, "onDestroy")
        stopForeground(STOP_FOREGROUND_REMOVE)
        runtime.start()
    }

    override fun onBind(intent: Intent) = null
}

internal fun Context.startAstralService() {
    val intent = Intent(this, MainService::class.java)
    if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
        startForegroundService(intent)
    } else {
        startService(intent)
    }
}

internal fun Context.stopAstralService() {
    val intent = Intent(this, MainService::class.java)
    stopService(intent)
}
