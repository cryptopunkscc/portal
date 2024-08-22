package cc.cryptopunks.portal

import android.app.Application
import cc.cryptopunks.portal.compose.composeModule
import cc.cryptopunks.portal.exception.exceptionModule
import cc.cryptopunks.portal.main.mainModule
import org.koin.android.ext.android.get
import org.koin.android.ext.koin.androidContext
import org.koin.core.context.startKoin
import java.lang.Thread.setDefaultUncaughtExceptionHandler

class AgentApp : Application() {

    override fun onCreate() {
        super.onCreate()
        startKoin {
            androidContext(applicationContext)
            modules(
                exceptionModule,
                mainModule,
                composeModule,
//                logcatModule,
//                nodeModule,
//                logModule,
//                configModule,
//                adminModule,
//                jsAppModule,
//                contactsModule,
//                warpdriveModule,
            )
        }
        setDefaultUncaughtExceptionHandler(get())
//        get<LogcatBackup>().start()
    }
}
