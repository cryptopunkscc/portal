package cc.cryptopunks.portal.app

import cc.cryptopunks.portal.Activities
import org.koin.dsl.bind
import org.koin.dsl.module

val applicationModule = module {
    single { ActivityStack() } bind Activities::class
}