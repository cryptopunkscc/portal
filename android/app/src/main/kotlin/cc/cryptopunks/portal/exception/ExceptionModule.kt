package cc.cryptopunks.portal.exception

import org.koin.core.module.dsl.factoryOf
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.bind
import org.koin.dsl.module

val exceptionModule = module {
    singleOf(::ExceptionsState)
    factoryOf(::ExceptionStorage)
    factoryOf(::ExceptionHandler).bind<Thread.UncaughtExceptionHandler>()
}
