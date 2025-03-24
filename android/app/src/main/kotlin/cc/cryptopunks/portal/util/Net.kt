package cc.cryptopunks.portal.util

import java.net.InterfaceAddress

val InterfaceAddress.CIDR: String
    get() = "${address.hostAddress.orEmpty().split("%")[0]}/$networkPrefixLength"

val Iterable<InterfaceAddress>.CIDRs: String
    get() = map(InterfaceAddress::CIDR).joinToString(separator = " ")