package cc.cryptopunks.portal.exception

import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.update

class ExceptionsState : (Throwable) -> Unit {

    private val errors = MutableStateFlow(emptyList<Throwable>())

    val current = errors.map { it.lastOrNull() }

    override fun invoke(e: Throwable) = plusAssign(e)

    operator fun plusAssign(e: Throwable) = errors.update { it + e }

    fun pop() = errors.update { it.dropLast(1) }
}
