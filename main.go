package main

import (
	"fmt"
	"os"
	"strings"
	"zonetree/cache"
	"zonetree/logger"
)

//var Log slog.Logger

var cfg cache.Config
var Log = logger.PrintDebugLog()

// var ZoneTree cmap.ConcurrentMap[string, Zone]
// var Zones cache.Map[cache.Zone]
var Zones = cache.NewZoneCache()

// var Cache = cmap.New[Server]()
// var Cache cache.Map[cache.Server]
var Cache = cache.NewServerCache()

//var Opt = logger.Options{IPv4only: true}

//var ZoneTree = cmap.New[Zone]()

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
