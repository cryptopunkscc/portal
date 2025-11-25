package core

import (
	"net"
	"strings"

	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	"github.com/cryptopunkscc/portal/api/mobile"
)

func interfaceAddrsFunc(api mobile.Api) func() ([]net.Addr, error) {
	return func() (out []net.Addr, err error) {
		cidrs, err := api.Net().Addresses()
		if err != nil {
			return
		}
		return parseCidrsToNetAddrs(cidrs)
	}
}

func netInterfacesFunc(api mobile.Api) func() ([]ether.NetInterface, error) {
	return func() (out []ether.NetInterface, err error) {
		i, err := api.Net().Interfaces()
		if err != nil {
			return
		}
		for n := i.Next(); n != nil; n = i.Next() {
			out = append(out, ether.NetInterface{
				Flags: net.Flags(n.Flags),
				Addrs: func() (addrs []net.Addr, err error) {
					return parseCidrsToNetAddrs(n.Addresses)
				},
			})
		}
		return
	}
}

func parseCidrsToNetAddrs(cidrs string) (addrs []net.Addr, err error) {
	for _, s := range strings.Fields(cidrs) {
		var ipNet *net.IPNet
		if _, ipNet, err = net.ParseCIDR(s); err != nil {
			return
		}
		addrs = append(addrs, ipNet)
	}
	return
}
