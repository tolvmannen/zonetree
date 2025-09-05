package api

import (
	//	"crypto/tls"
	//	"os"

	//	"fmt"
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

	router.GET("/cache/list/*zone", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		zone := strings.TrimLeft(c.Param("zone"), "/")
		if zone != "" {
			zone = cache.MakeFQDN(zone)
		}

		// default outstr if nothing returned from cache
		outstr := "Zone not in cache:[" + zone + "]\n"

		if z, ok := Zones.Get(zone); ok {
			outstr, _ = z.ToPrettyJson()
		}

		c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

	})

	router.GET("/test/*zone", func(c *gin.Context) {
		// trim any leading slash (applies when no 'name' is provided)
		zone := strings.TrimLeft(c.Param("zone"), "/")

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
