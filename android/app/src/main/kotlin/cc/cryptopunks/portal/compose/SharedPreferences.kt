package cc.cryptopunks.portal.compose

import android.content.SharedPreferences
import androidx.compose.runtime.MutableState
import androidx.core.content.edit

fun SharedPreferences.mutableStateOf(
    name: String,
    default: Boolean
) = mutableStateOf(
    default = default,
    put = { putBoolean(name, it) },
    get = { getBoolean(name, default) }
)

fun <V> SharedPreferences.mutableStateOf(
    default: V,
    put: SharedPreferences.Editor.(V) -> Unit,
    get: SharedPreferences.() -> V,
) = SharedPreferencesMutableState(
    state = androidx.compose.runtime.mutableStateOf(default),
    sharedPreferences = this,
    put = put,
    get = get
)

class SharedPreferencesMutableState<V>(
    private val state: MutableState<V>,
    private val sharedPreferences: SharedPreferences,
    private val put: SharedPreferences.Editor.(V) -> Unit,
    get: SharedPreferences.() -> V,
) : MutableState<V> {

    override var value: V
        get() = state.value
        set(value) {
            state.value = value
            sharedPreferences.edit {
                put(this, value)
            }
        }

    init {
        value = sharedPreferences.get()
    }

    override fun component1(): V = value

    override fun component2(): (V) -> Unit = { value = it }
}
