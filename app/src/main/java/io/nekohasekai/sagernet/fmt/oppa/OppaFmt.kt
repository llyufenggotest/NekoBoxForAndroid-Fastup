package io.nekohasekai.sagernet.fmt.oppa

import io.nekohasekai.sagernet.ktx.linkBuilder
import io.nekohasekai.sagernet.database.DataStore
import moe.matsuri.nb4a.SingBoxOptions
import io.nekohasekai.sagernet.ktx.toLink
import io.nekohasekai.sagernet.ktx.urlSafe
import okhttp3.HttpUrl.Companion.toHttpUrlOrNull
import org.json.JSONObject
import java.util.Base64

data class OppaProvider(
    val name: String,
    val apiToken: String,
    val apiEndpoint: String,
    val encryptionKey: String,
)

fun parseOppaProvider(link: String): OppaProvider {
    require(link.startsWith("oppa://")) { "invalid Oppa provider link" }
    val payload = link.removePrefix("oppa://")
    require(payload.isNotBlank() && !payload.contains('@')) { "invalid Oppa provider payload" }
    val decoded = runCatching {
        Base64.getUrlDecoder().decode(payload.padEnd((payload.length + 3) / 4 * 4, '='))
    }.getOrElse { throw IllegalArgumentException("invalid Oppa provider encoding", it) }
    val json = runCatching { JSONObject(String(decoded, Charsets.UTF_8)) }
        .getOrElse { throw IllegalArgumentException("invalid Oppa provider JSON", it) }
    val name = json.optString("name").trim()
    val apiToken = json.optString("apiToken").trim()
    val apiEndpoint = json.optString("apiEndPoint").trim().removePrefix("https://").trimEnd('/')
    val encryptionKey = json.optString("encryptionKey").trim()
    require(name.isNotEmpty()) { "missing Oppa provider name" }
    require(apiToken.isNotEmpty()) { "missing Oppa provider token" }
    require(apiEndpoint.isNotEmpty()) { "missing Oppa provider endpoint" }
    require(encryptionKey.isNotEmpty()) { "missing Oppa provider encryption key" }
    return OppaProvider(name, apiToken, apiEndpoint, encryptionKey)
}

fun parseOppaNode(url: String): OppaBean {
    require(url.startsWith("oppa://")) { "invalid Oppa node link" }
    val link = url.replaceFirst("oppa://", "https://").toHttpUrlOrNull()
        ?: throw IllegalArgumentException("invalid Oppa node link")
    val password = link.username
    require(password.isNotEmpty() && password.toByteArray(Charsets.UTF_8).size <= 4096) { "Oppa password must contain 1–4096 UTF-8 bytes" }
    return OppaBean().apply {
        serverAddress = link.host
        serverPort = link.port
        this.password = password
        preConnect = link.queryParameter("preconnect")?.toIntOrNull()?.coerceIn(0, 64) ?: 8
        sni = link.queryParameter("sni") ?: ""
        allowInsecure = link.queryParameter("insecure")?.let { it == "1" || it.equals("true", true) } ?: false
        certificatePinSHA256 = link.queryParameter("pin_sha256") ?: ""
        name = link.fragment ?: ""
        initializeDefaultValues()
    }
}

fun buildSingBoxOutboundOppaBean(bean: OppaBean): SingBoxOptions.Outbound_OppaOptions {
    return SingBoxOptions.Outbound_OppaOptions().apply {
        type = "oppa"
        server = bean.serverAddress
        server_port = bean.serverPort
        password = bean.password
        pre_connect = bean.preConnect
        pin_cert_sha256 = bean.certificatePinSHA256.ifBlank { null }
        tls = SingBoxOptions.OutboundTLSOptions().apply {
            enabled = true
            server_name = bean.sni.ifBlank { null }
            insecure = bean.allowInsecure || DataStore.globalAllowInsecure
            alpn = null
        }
    }
}

fun OppaBean.toUri(): String {
    require(password.isNotEmpty() && password.toByteArray(Charsets.UTF_8).size <= 4096) { "Oppa password must contain 1–4096 UTF-8 bytes" }
    val builder = linkBuilder()
        .host(serverAddress)
        .port(serverPort)
        .username(password)
    if (preConnect != 8) builder.addQueryParameter("preconnect", preConnect.toString())
    if (sni.isNotBlank()) builder.addQueryParameter("sni", sni)
    if (allowInsecure) builder.addQueryParameter("insecure", "1")
    if (certificatePinSHA256.isNotBlank()) builder.addQueryParameter("pin_sha256", certificatePinSHA256)
    if (name.isNotBlank()) builder.encodedFragment(name.urlSafe())
    return builder.toLink("oppa")
}
