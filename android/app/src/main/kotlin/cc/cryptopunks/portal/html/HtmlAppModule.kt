package cc.cryptopunks.portal.html

import cc.cryptopunks.portal.StartHtmlApp
import org.koin.core.module.dsl.singleOf
import org.koin.dsl.binds
import org.koin.dsl.module

val htmlAppModule = module {
    singleOf(::HtmlAppRepository) binds arrayOf(
        StartHtmlApp::class,
    )
}