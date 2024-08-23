package cc.cryptopunks.portal.main

import cc.cryptopunks.portal.CoreEvents
import cc.cryptopunks.portal.core.mobile.Event
import kotlinx.coroutines.flow.MutableStateFlow

class MainEvents : CoreEvents,
    MutableStateFlow<Event> by MutableStateFlow(Event())