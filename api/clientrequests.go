package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"zonetree/cache"
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

	if tr.Undelegated {
		text = UndelegatedTest(tr)
	} else {
		text = SmallTest(tr)
	}

	//text := tr.String()
	elapsed := time.Since(start)

	//outstr = "Tested Zone:[" + zone + "] (" + elapsed.String() + ")\n"
	outstr = text + " \n\nTook: (" + elapsed.String() + ")\n"

	c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

}

func UndelegatedTest(tr zonetests.TestRequest) string {
	// If the Undelegated flag is set, there is no need to
	// traverse the DNS tree.
	for _, z := range tr.Zones {
		cache.PrepUndelegatedZone(z.Name, &cfg, z.UNS)
	}

	return "Undelegated test: "
}

func SmallTest(tr zonetests.TestRequest) string {

	var text string

	for _, z := range tr.Zones {

		cache.BuildZoneCache(z.Name, &cfg)

		//t := zonetests.NewRunner(&cfg, zone.Name)
		for _, testname := range tr.Tests {
			//t.AddTest(test)
			text += "Zone: " + z.Name + " - " + testname + "\n"
		}

	}
	return text
}
