package io.nekohasekai.sagernet.fmt.trojan

import io.nekohasekai.sagernet.fmt.v2ray.toUriVMessVLESSTrojan
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class TrojanFastupFmtTest {

    @Test
    fun fastupSuffixSurvivesImportAndExport() {
        val link = "trojan://synthetic-password%23fastup@node.example:443?security=tls&type=tcp#Fastup"
        val bean = parseTrojan(link)

        assertEquals("synthetic-password#fastup", bean.password)
        bean.initializeDefaultValues()
        val exported = bean.toUriVMessVLESSTrojan(true)
        assertTrue(exported.contains("synthetic-password%23fastup@"))
        assertEquals(bean.password, parseTrojan(exported).password)
    }

    @Test
    fun standardTrojanPasswordIsUnchanged() {
        val bean = parseTrojan("trojan://ordinary-password@node.example:443?security=tls&type=tcp#Standard")
        assertEquals("ordinary-password", bean.password)
    }
}
