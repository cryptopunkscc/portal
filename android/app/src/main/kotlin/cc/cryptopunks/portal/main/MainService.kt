package cc.cryptopunks.portal.main

import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.util.Log
import cc.cryptopunks.portal.CoreEvents
import cc.cryptopunks.portal.core.mobile.Event
import cc.cryptopunks.portal.core.mobile.Mobile
import cc.cryptopunks.portal.core.mobile.Core
import cc.cryptopunks.portal.exception.ExceptionsState
import cc.cryptopunks.portal.util.acquireMulticastWakeLock
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.cancel
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.mapNotNull
import kotlinx.coroutines.launch
import org.koin.android.ext.android.inject

class MainService : Service() {

    private val scope = CoroutineScope(Dispatchers.IO)
    private val tag = javaClass.simpleName
    private val runtime: Core by inject()
    private val events: CoreEvents by inject()
    private val wakeLock by lazy { acquireMulticastWakeLock("astral") }
    private val exceptions: ExceptionsState by inject()

    override fun onCreate() {
        Log.d(tag, "onCreate")
        startForegroundNotification()
        wakeLock.acquire()
        scope.launch {
            events.collect { event -> Log.d(tag, "event: ${event.msg}", event.err) }
        }
        scope.launch {
            events.mapNotNull(Event::getErr).collect(exceptions)
        }
        scope.launch {
            events.first { event -> event.msg == Mobile.STOPPED }
            stopSelf()
        }
        scope.launch {
            runtime.start()
        }
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        Log.d(tag, "onStartCommand")
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
        runtime.stop()
        wakeLock.release()
        scope.cancel()
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
