package cc.cryptopunks.portal.main

import android.app.Activity
import android.app.Application.ActivityLifecycleCallbacks
import android.os.Bundle
import java.util.WeakHashMap

class ActivityStack : GetCurrentActivity, ActivityLifecycleCallbacks {

    private val stack = WeakHashMap<Activity, Entry>()

    override fun invoke(): Activity? = stack.toMap()
        .asSequence().groupBy { it.value.state }
        .asSequence().sortedBy { it.key }
        .flatMap { (_, entries) -> entries.sortedBy { it.value.time } }
        .toList().lastOrNull()?.key

    override fun onActivityCreated(activity: Activity, savedInstanceState: Bundle?) {
        stack[activity] = Entry()
    }

    override fun onActivityStarted(activity: Activity) = activity set State.Started
    override fun onActivityResumed(activity: Activity) = activity set State.Resumed
    override fun onActivityPaused(activity: Activity) = activity set State.Started
    override fun onActivityStopped(activity: Activity) = activity set State.Created
    override fun onActivitySaveInstanceState(activity: Activity, outState: Bundle) = Unit
    override fun onActivityDestroyed(activity: Activity) {
        stack.remove(activity)
    }

    private infix fun Activity.set(state: State) {
        stack.getOrDefault(this, Entry()).state = state
    }

    private class Entry {
        var time: Long = System.currentTimeMillis(); private set
        var state: State = State.Created
            set(value) {
                field = value
                time = System.currentTimeMillis()
            }
    }
    private enum class State { Created, Started, Resumed }
}