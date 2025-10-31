package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"zonetree/zonetests"
)

func RunTest(c *gin.Context) {
	// Dirty execution time check
	start := time.Now()

	//var zt zonetests.TestRunner
	var tr zonetests.TestRequest
	var outstr string
	if err := c.ShouldBindJSON(&tr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//cache.BuildZoneCache(zone, &cfg)

	var text string

	for _, zone := range tr.Zones {

		text += "Zone: " + zone.Name + "\n"
		//t := zonetests.NewRunner(&cfg, zone.Name)
		for _, testname := range tr.Tests {
			//t.AddTest(test)
			text += "  - " + testname + "\n"
		}

	}

	//text := tr.String()
	elapsed := time.Since(start)

	//outstr = "Tested Zone:[" + zone + "] (" + elapsed.String() + ")\n"
	outstr = text + " /n/n(" + elapsed.String() + ")\n"

	c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

}
