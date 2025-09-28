package api

import (
	//	"crypto/tls"
	//	"os"

	//	"fmt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strings"
	//	"time"

	//	"github.com/gin-contrib/cors"
	//	"github.com/gin-contrib/static"
	//	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/acme/autocert"
	// "gopkg.in/yaml.v3"
	"zonetree/cache"
	"zonetree/html"
	"zonetree/logger"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

var cfg cache.Config
var Log = logger.PrintDebugLog()
var Zones = cache.NewZoneCache()
var Cache = cache.NewServerCache()

func Run() {

	cfg = cache.Init(&Log, Zones, Cache)

	router := gin.Default()

	router.GET("/conf/show", func(c *gin.Context) {

		outstr := cfg.RunningConf()

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/conf/load/*file", func(c *gin.Context) {

		file := strings.TrimLeft(c.Param("file"), "/")

		var outstr string

		err := cfg.Load(file)
		if err != nil {
			outstr = err.Error()
		} else {
			outstr = "loaded conf file" + file
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/cache/tree/*zone", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		zone := strings.TrimLeft(c.Param("zone"), "/")
		if zone != "" {
			zone = cache.ToFQDN(strings.ToLower(zone))
		}

		// default outstr if nothing returned from cache
		//outstr := "Zone not in cache:[" + zone + "]\n"
		list := cache.DigPath(zone)
		tree := cfg.ZoneCutPath(list)
		tree = append([]string{"."}, tree...)

		fmt.Printf("\n%v\n", tree)

		// Recursively (o_O) go through the list
		var tr func(l []string) html.Node

		tr = func(l []string) html.Node {

			var HN html.Node

			// Get Current zone
			var qns bool
			if z, ok := Zones.Get(l[0]); ok {

				HN.Name = z.Name

				qns = false
				for _, ns := range z.NSIP {
					//fmt.Printf("qns: %v : %v\n", qns, ns.ZoneStatus)
					if qns == false && (ns.ZoneStatus == 200 || ns.ZoneStatus == 0) && len(l) > 1 {
						n := tr(l[1:])
						HN.Children = append(HN.Children, n)
						qns = true

					} else {
						n := html.Node{Name: ns.Name + " (" + ns.IP + ")", Parent: &HN}
						HN.Children = append(HN.Children, n)
					}
				}

			}

			return HN
		}

		var nodetree html.Node = tr(tree)

		outstr, err := json.MarshalIndent(nodetree, "", "  ")
		if err != nil {
			outstr = []byte(err.Error())
		}
		/*
			if z, ok := Zones.Get(zone); ok {
				list := cache.DigPath(zone)
				tree := cfg.ZoneCutPath(list)
				outstr = "Zone cuts for" + z.Name + ":\n"
				outstr += strings.Join(tree, "-->") + "\n"
			}

		*/

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/cache/list/*zone", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		zone := strings.TrimLeft(c.Param("zone"), "/")
		if zone != "" {
			zone = cache.ToFQDN(strings.ToLower(zone))
		}

		// default outstr if nothing returned from cache
		outstr := "Zone not in cache:[" + zone + "]\n"

		if z, ok := Zones.Get(zone); ok {
			outstr, _ = z.ToPrettyJson()
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/cache/clear/*zone", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		zone := strings.TrimLeft(c.Param("zone"), "/")
		if zone != "" {
			zone = cache.ToFQDN(strings.ToLower(zone))
		}

		// default outstr if nothing returned from cache
		outstr := "Zone not in cache:[" + zone + "]\n"

		if z, ok := Zones.Get(zone); ok {
			if z.Name != "." {
				// Dont delete the ROOT
				Zones.Remove(zone)
				outstr = "Zone [" + z.Name + "] removed from cache"
			} else {
				outstr = "ERROR: ROOT cannot be deleted. Restart to load new data from root-hints."
			}
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/cache/reset", func(c *gin.Context) {

		var outstr string

		for z := range cfg.Zones.IterBuffered() {
			if zone, ok := Zones.Get(z.Value.Name); ok {
				// Dont delete the ROOT
				if zone.Name != "." {
					Zones.Remove(zone.Name)
					outstr += "Zone [" + zone.Name + "] removed from cache\n"
				}
			}
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/cache/dump", func(c *gin.Context) {

		var outstr string

		for z := range cfg.Zones.IterBuffered() {
			if zone, ok := Zones.Get(z.Value.Name); ok {
				jstr, _ := zone.ToPrettyJson()
				outstr += jstr
			}
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/test/*zone", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		zone := strings.ToLower(strings.TrimLeft(c.Param("zone"), "/"))

		cache.BuildZoneCache(zone, &cfg)

		outstr := "Testing Zone:[" + zone + "]\n"
		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	server := &http.Server{
		Addr:    ":7777",
		Handler: router,
	}
	log.Fatal(server.ListenAndServe())

}
