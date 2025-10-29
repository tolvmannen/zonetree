package zonetests

import (
	"fmt"
	"net/netip"
	"zonetree/cache"
)

func (ts *TestSuite) Address01(ipaddr string) {

	var tbl []cache.SuIP

	files := []string{
		"assets/iana-ipv6-special-registry-1.csv",
		"assets/iana-ipv4-special-registry-1.csv",
	}

	for _, file := range files {
		tbl = append(cache.ReadCSV(file), tbl...)
	}

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
