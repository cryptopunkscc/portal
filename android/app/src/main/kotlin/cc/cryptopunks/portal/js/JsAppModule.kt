package cc.cryptopunks.portal.js

import org.koin.core.module.dsl.factoryOf
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.module

val jsAppModule = module {
    singleOf(::JsAppsManager)
    factoryOf(::JsAppActivities)
}
