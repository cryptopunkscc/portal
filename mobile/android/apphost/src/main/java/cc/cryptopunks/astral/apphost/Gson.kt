package cc.cryptopunks.astral.apphost

import com.google.gson.Gson
import com.google.gson.reflect.TypeToken

class GsonAppHostClient(
    client: AppHostClient,
    gson: Gson = defaultGson,
) : AppHostClientAdapter<ConnSerializer>(client) {
    private val serializer = GsonSerializer(gson)
    override fun convert(conn: Conn) = conn.serializer(serializer)
}

fun Conn.gsonSerializer() = serializer(gsonSerializer)

val gsonSerializer: Serializer by lazy {
    GsonSerializer(defaultGson)
}

val defaultGson by lazy {
    Gson().newBuilder()
        .setFieldNamingStrategy { it.name.replaceFirstChar(Char::uppercase) }
        .create()
}

private class GsonSerializer(
    private val gson: Gson = defaultGson,
) : Serializer {
    override fun encode(any: Any): String = gson.toJson(any)
    override fun <T> decode(string: String, type: Class<T>): T = gson.fromJson(string, type)

    override fun <T> decodeList(string: String, type: Class<T>): List<T> = when {
        string.isBlank() -> emptyList()
        else -> gson.fromJson<Array<T>>(
            string, TypeToken.getArray(
                TypeToken.get(type).type
            ).type
        )?.toList() ?: emptyList()
    }

    override fun <K, V> decodeMap(string: String, key: Class<K>, value: Class<V>): Map<K, V> =
        gson.fromJson(
            string, TypeToken.getParameterized(
                TypeToken.get(Map::class.java).type, key, value
            ).type
        )
}
