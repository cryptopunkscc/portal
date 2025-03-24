package cc.cryptopunks.portal.agent

import cc.cryptopunks.portal.core.mobile.Net
import cc.cryptopunks.portal.core.mobile.NetInterface
import cc.cryptopunks.portal.core.mobile.NetInterfaceIterator
import cc.cryptopunks.portal.util.CIDRs
import java.net.NetworkInterface

class AgentNet : Net {

    override fun addresses(): String {
        return NetworkInterface.getNetworkInterfaces().toList()
            .flatMap(NetworkInterface::getInterfaceAddresses)
            .CIDRs
    }

    override fun interfaces(): NetInterfaceIterator {
        val iterator = NetworkInterface.getNetworkInterfaces().asSequence()
            .map(NetworkInterface::toNetworkInterfaceInfo)
            .onEach { println(it) }
            .iterator()
        return NetInterfaceIterator {
            if (iterator.hasNext()) iterator.next() else null
        }
    }
}

private fun NetworkInterface.toNetworkInterfaceInfo() = NetInterface().also { i ->
    i.index = index.toLong()
    i.mtu = mtu.toLong()
    i.name = name
    i.hardwareAddr = hardwareAddress
    i.addresses = interfaceAddresses.CIDRs
    i.flags = mapOf(
        1 to isUp,                  // net.FlagUp
        2 to hasBroadcast,          // net.FlagBroadcast
        4 to isLoopback,            // net.FlagLoopback
        8 to isPointToPoint,        // net.FlagPointToPoint
        16 to supportsMulticast(),  // net.FlagMulticast
    ).filterValues(true::and).keys.fold(0, Int::or).toLong()
}

private val NetworkInterface.hasBroadcast get() = interfaceAddresses.any { it.broadcast != null }
