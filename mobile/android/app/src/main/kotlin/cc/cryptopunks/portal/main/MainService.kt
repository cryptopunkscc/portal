package cc.cryptopunks.portal.main

import android.app.Service
import android.content.Context
import android.content.Intent
import android.os.Build
import android.util.Log
import cc.cryptopunks.portal.Errors
import cc.cryptopunks.portal.Status
import cc.cryptopunks.portal.core.mobile.Mobile
import cc.cryptopunks.portal.core.mobile.Core
import cc.cryptopunks.portal.exception.ExceptionsState
import cc.cryptopunks.portal.util.acquireMulticastWakeLock
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.cancel
import kotlinx.coroutines.flow.filter
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.timeout
import kotlinx.coroutines.launch
import org.koin.android.ext.android.inject
import kotlin.time.Duration.Companion.seconds

class MainService : Service() {

    private val scope = CoroutineScope(Dispatchers.IO)
    private val tag = javaClass.simpleName
    private val runtime: Core by inject()
    private val status: Status by inject()
    private val errors: Errors by inject()
    private val wakeLock by lazy { acquireMulticastWakeLock("astral") }
    private val exceptions: ExceptionsState by inject()

    override fun onCreate() {
        Log.d(tag, "onCreate")
        startForegroundNotification()
        wakeLock.acquire()
        scope.launch {
            status.collect { Log.d(tag, "status: $it") }
        }
        scope.launch {
            errors.map(::Exception).collect(exceptions)
        }
        scope.launch {
            runCatching {
                status.filter { it == Mobile.STARTING }.timeout(3.seconds).first()
            }.onSuccess {
                status.first { it == Mobile.STOPPED }
            }
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
