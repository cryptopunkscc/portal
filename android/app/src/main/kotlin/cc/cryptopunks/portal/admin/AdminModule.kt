package cc.cryptopunks.portal.admin

import org.koin.core.module.dsl.singleOf
import org.koin.dsl.module

val adminModule = module {
    singleOf(::AdminClient)
    singleOf(::AdminPreferences)
    singleOf(::AdminViewModel)
}
