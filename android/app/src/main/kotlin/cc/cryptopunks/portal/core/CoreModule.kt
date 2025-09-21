package cc.cryptopunks.portal.core

import cc.cryptopunks.portal.Errors
import cc.cryptopunks.portal.Status
import cc.cryptopunks.portal.core.core.Core
import cc.cryptopunks.portal.core.mobile.Api
import org.koin.core.module.dsl.factoryOf
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.bind
import org.koin.dsl.module

val coreModule = module {
    factoryOf(::CoreApi) bind Api::class
    singleOf(::CoreStatus) bind Status::class
    singleOf(::CoreErrors) bind Errors::class
    singleOf(Core::create)
}
