package cc.cryptopunks.portal.compose

import cc.cryptopunks.portal.exception.ErrorsScreen
import org.koin.dsl.module

val composeModule = module {
    single {
        ComposeApi(
            Theme = { AstralTheme(it) },
            Errors = { ErrorsScreen() },
//            Contacts = { ContactsScreen(select = it) }
        )
    }
}
