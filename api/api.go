package api

import (
	"log"
	"net/http"
	//"time"

	"strings"

	//	"github.com/gin-contrib/cors"
	//	"github.com/gin-contrib/static"
	//	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	// "golang.org/x/crypto/acme/autocert"
	// "gopkg.in/yaml.v3"
	"zonetree/cache"
	//"zonetree/logger"
	//"zonetree/zonetests"
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

	router.GET("/conf/show", ShowConf)

	router.GET("/conf/load/*file", LoadConf)

	router.POST("/run", RunTest)

	router.GET("/cache/list/*zone", PrintJzone)

	router.GET("/cache/clear/*zone", ClearZone)

	router.GET("/cache/reset", ClearCache)

	router.GET("/cache/dump", DumpCache)

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
