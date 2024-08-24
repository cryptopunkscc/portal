package cc.cryptopunks.portal.main

import cc.cryptopunks.portal.core.factory.Factory
import cc.cryptopunks.portal.core.mobile.Api
import org.koin.core.module.dsl.factoryOf
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.bind
import org.koin.dsl.module

val mainModule = module {
    factoryOf(Factory::runtime)
    factoryOf(::MainApi) bind Api::class
    singleOf(::MainEvents)
    singleOf(::MainPermissions)
}
