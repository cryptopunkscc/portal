package cc.cryptopunks.portal.api

import kotlinx.coroutines.flow.Flow

fun interface ServiceStatus : (String) -> Flow<Boolean>
