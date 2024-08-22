package cc.cryptopunks.portal.exception

class ExceptionHandler(
    private val state: ExceptionsState,
    private val storage: ExceptionStorage,
) : Thread.UncaughtExceptionHandler {

    override fun uncaughtException(t: Thread, e: Throwable) {
        e.printStackTrace()
        state += e
        storage.save(e)
    }
}
