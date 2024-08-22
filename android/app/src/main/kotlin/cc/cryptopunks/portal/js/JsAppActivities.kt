package cc.cryptopunks.portal.js

import android.content.Context
import androidx.core.content.edit
import cc.cryptopunks.portal.js.JsAppActivity

class JsAppActivities(context: Context) {

    private val preferences = context.getSharedPreferences("jsAppActivities", Context.MODE_PRIVATE)

    operator fun get(appId: String): Int = preferences.getInt(appId, -1)

    operator fun set(appId: String, activityId: Int) = preferences.edit {
        if (activityId < 0) remove(appId)
        else putInt(appId, activityId)
    }

    fun nextId(): Int {
        val used = preferences.all.mapNotNull { it.value as? Int }.sorted()
        var i = -1
        for (id in used) if (id != ++i) return i
        if (i == JsAppActivity.Limit) throw AllSlotsReserved()
        return i
    }

    class AllSlotsReserved : Exception()
}
