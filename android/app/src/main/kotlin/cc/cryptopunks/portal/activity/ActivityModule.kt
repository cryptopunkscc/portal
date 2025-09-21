package cc.cryptopunks.portal.activity

import cc.cryptopunks.portal.Activities
import org.koin.dsl.bind
import org.koin.dsl.module

val activityModule = module {
    single { ActivityStack() } bind Activities::class
}