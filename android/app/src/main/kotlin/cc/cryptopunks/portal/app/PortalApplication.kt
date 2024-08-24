package cc.cryptopunks.portal.app

import android.app.Application
import cc.cryptopunks.portal.compose.composeModule
import cc.cryptopunks.portal.core.factory.Factory
import cc.cryptopunks.portal.core.mobile.Runtime
import cc.cryptopunks.portal.exception.exceptionModule
import cc.cryptopunks.portal.html.htmlAppModule
import cc.cryptopunks.portal.main.mainModule
import org.koin.android.ext.android.get
import org.koin.android.ext.koin.androidContext
import org.koin.core.context.startKoin
import java.lang.Thread.setDefaultUncaughtExceptionHandler

class PortalApplication : Application() {

    override fun onCreate() {
        super.onCreate()
        startKoin {
            androidContext(applicationContext)
            modules(
                exceptionModule,
                composeModule,
                mainModule,
                htmlAppModule,
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
        get<Runtime>().install()
//        get<LogcatBackup>().start()
    }
}
