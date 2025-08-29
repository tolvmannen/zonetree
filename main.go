package main

import (
	"fmt"
	"os"
	"strings"
	"zonetree/cache"
	"zonetree/logger"
)

var cfg cache.Config
var Log = logger.PrintDebugLog()
var Zones = cache.NewZoneCache()
var Cache = cache.NewServerCache()

func main() {

	var domain string

	if len(os.Args) > 1 {
		//domain = dns.Fqdn(os.Args[1])
		domain = cache.MakeFQDN(os.Args[1])
		if strings.ToUpper(domain) == "ROOT." {
			domain = "." // use root if no domain given
		}

	} else {
		domain = "." // use root if no domain given
	}

	//ZoneTree, Cache, Log = Init()
	cfg = cache.Init(&Log, Zones, Cache)

	cache.BuildZoneCache(domain, &cfg)

	/*
		fmt.Printf("\n\n-----IN CACHE ------\n\n")
		for m := range Cache.IterBuffered() {
			fmt.Printf("%v\n", m)
		}
	*/

	z, ok := Zones.Get(domain)
	if ok {
		//z.Print()
		j, _ := z.ToPrettyJson()
		fmt.Printf("%v\n", j)

	}

	/*
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	*/
	//zone.Print()

}
