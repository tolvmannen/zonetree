package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"zonetree/dig"
)

// ShowConf
//
// Prints the currently loaded config
func ShowConf(c *gin.Context) {
	outstr := cfg.RunningConf()
	c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))
}

// LoadConf
//
// (Re)loads config from file
func LoadConf(c *gin.Context) {
	file := strings.TrimLeft(c.Param("file"), "/")

	var outstr string

	err := cfg.Load(file)
	if err != nil {
		outstr = err.Error()
	} else {
		outstr = "loaded conf file" + file
	}

	c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))
}

// PrintJzone
//
// Pretty prints the Zone object in JSON format
func PrintJzone(c *gin.Context) {
	// trim any leading slash (applies when no 'name' is provided)
	zone := strings.TrimLeft(c.Param("zone"), "/")
	if zone != "" {
		zone = dig.ToFQDN(strings.ToLower(zone))
	}

	// default outstr if nothing returned from cache
	outstr := "Zone not in cache:[" + zone + "]\n"

	if z, ok := cfg.Zones.Get(zone); ok {
		outstr, _ = z.ToPrettyJson()
	}

	c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))

}

// ClearZone
//
// Delete zone from cache
func ClearZone(c *gin.Context) {
	// trim any leading slash (applies when no 'name' is provided)
	zone := strings.TrimLeft(c.Param("zone"), "/")
	if zone != "" {
		zone = dig.ToFQDN(strings.ToLower(zone))
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
}

// ClearCache
//
// Delete all zones (exept ROOT) from cache
func ClearCache(c *gin.Context) {
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
}

// DumpCache
//
// Pretty print all zones in cache
func DumpCache(c *gin.Context) {
	var outstr string

	for z := range cfg.Zones.IterBuffered() {
		if zone, ok := cfg.Zones.Get(z.Value.Name); ok {
			jstr, _ := zone.ToPrettyJson()
			outstr += jstr
		}
	}

	c.Data(http.StatusOK, ContentTypeHTML, []byte(outstr))
}
