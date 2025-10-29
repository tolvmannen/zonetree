package api

import (
	//	"crypto/tls"
	//	"os"

	//	"fmt"
	//	"encoding/json"
	"log"
	"net/http"
	"time"

	"strings"

	//	"github.com/gin-contrib/cors"
	//	"github.com/gin-contrib/static"
	//	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/acme/autocert"
	// "gopkg.in/yaml.v3"
	"zonetree/cache"
	//"zonetree/logger"
	"zonetree/zonetests"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

var cfg cache.Config

func Run() {

	cfg = cache.Init()

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

	router.POST("/ut/*test", func(c *gin.Context) {
		// Dirty execution time check

		test := strings.TrimLeft(c.Param("test"), "/")

		var zt zonetests.TestSuite
		var outstr string
		if err := c.ShouldBindJSON(&zt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		/*

			zone := strings.ToLower(zt.Zone)
			cache.BuildZoneCache(zone, &cfg)

			elapsed := time.Since(start)

			outstr = "Tested Zone:[" + zone + "] (" + elapsed.String() + ")\n"
			//outstr = "Test zone:" + zt.Zone + "\n"

		*/

		var iplist = []string{
			"127.0.0.1",
			"203.0.113.15",
			"192.0.0.8",
			"192.0.1.1",
			"192.88.99.2",
			"192.88.99.9",
			"192.0.0.8",
			"213.50.29.170",
		}

		for _, ipaddr := range iplist {
			zt.Address01(ipaddr, cfg.SuIP)
		}

		outstr = "Doing test:[" + test + "] )\n"

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

		//c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))
	})

	router.POST("/run", func(c *gin.Context) {
		// Dirty execution time check
		start := time.Now()

		var zt zonetests.TestSuite
		var outstr string
		if err := c.ShouldBindJSON(&zt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		zone := strings.ToLower(zt.Zone)

		cache.BuildZoneCache(zone, &cfg)

		elapsed := time.Since(start)

		outstr = "Tested Zone:[" + zone + "] (" + elapsed.String() + ")\n"
		//outstr = "Test zone:" + zt.Zone + "\n"

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

		//c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))
	})

	router.GET("/cache/list/*zone", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		zone := strings.TrimLeft(c.Param("zone"), "/")
		if zone != "" {
			zone = cache.ToFQDN(strings.ToLower(zone))
		}

		// default outstr if nothing returned from cache
		outstr := "Zone not in cache:[" + zone + "]\n"

		if z, ok := cfg.Zones.Get(zone); ok {
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

		if z, ok := cfg.Zones.Get(zone); ok {
			if z.Name != "." {
				// Dont delete the ROOT
				cfg.Zones.Remove(zone)
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
			if zone, ok := cfg.Zones.Get(z.Value.Name); ok {
				// Dont delete the ROOT
				if zone.Name != "." {
					cfg.Zones.Remove(zone.Name)
					outstr += "Zone [" + zone.Name + "] removed from cache\n"
				}
			}
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/cache/dump", func(c *gin.Context) {

		var outstr string

		for z := range cfg.Zones.IterBuffered() {
			if zone, ok := cfg.Zones.Get(z.Value.Name); ok {
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
