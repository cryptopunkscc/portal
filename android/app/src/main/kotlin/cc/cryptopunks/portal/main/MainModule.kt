package cc.cryptopunks.portal.main

import cc.cryptopunks.portal.core.factory.Factory
import org.koin.core.module.dsl.factoryOf
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.bind
import org.koin.dsl.module

val mainModule = module {
    factoryOf(Factory::runtime)
    singleOf(::MainPermissions)
    singleOf(::MainApi)
    singleOf(::ActivityStack).bind<GetCurrentActivity>()

}
