package cc.cryptopunks.portal.main

import cc.cryptopunks.portal.CoreEvents
import cc.cryptopunks.portal.agent.AgentApi
import cc.cryptopunks.portal.core.core.Core
import cc.cryptopunks.portal.core.mobile.Api
import org.koin.core.module.dsl.factoryOf
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.bind
import org.koin.dsl.module

val mainModule = module {
    factoryOf(Core::create)
    factoryOf(::AgentApi) bind Api::class
    singleOf(::MainEvents) bind CoreEvents::class
    singleOf(::MainPermissions)
}
