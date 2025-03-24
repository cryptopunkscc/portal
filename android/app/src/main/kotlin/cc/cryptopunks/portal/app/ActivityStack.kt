package cc.cryptopunks.portal.app

import android.app.Activity
import android.app.Application
import android.os.Bundle
import cc.cryptopunks.portal.Activities

class ActivityStack(
    private val stack: MutableList<Activity> = mutableListOf()
) : Application.ActivityLifecycleCallbacks,
    Activities,
    List<Activity> by stack {

    override fun onActivityCreated(activity: Activity, savedInstanceState: Bundle?) {
        stack += activity
    }

    override fun onActivityDestroyed(activity: Activity) {
        stack -= activity
    }

    override fun onActivityStarted(activity: Activity) = Unit
    override fun onActivityResumed(activity: Activity) = Unit
    override fun onActivityPaused(activity: Activity) = Unit
    override fun onActivityStopped(activity: Activity) = Unit
    override fun onActivitySaveInstanceState(activity: Activity, outState: Bundle) = Unit
}