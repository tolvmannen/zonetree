package zonetests

import (
	"fmt"
	"net/netip"
	"zonetree/cache"
)

func (zt *ZoneTest) Address01() {

}

func (ts TestRunner) DepAddress01(ipaddr string, tbl []cache.SuIP) {

	for _, r := range tbl {

		net, err := netip.ParsePrefix(r.Block)

		if err != nil {
			panic(err)
		}

		ip, err := netip.ParseAddr(ipaddr)

		if err != nil {
			panic(err)
		}

		b := net.Contains(ip)
		if b == true {
			fmt.Printf("IP: %v - Net: %v - Routable? %v\n", ip, net, r.Global)
		}

	}
	//var t TestResult

}
