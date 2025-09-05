package cache

import (
	"fmt"
	"math/rand"
	"strconv"
	"zonetree/logger"
)

type Config struct {
	Log          logger.Logger
	Zones        Map[Zone]
	Cache        Map[Server]
	IPv4only     bool     `json:"IPv4only"`
	IPv6only     bool     `json:"IPv6only"`
	ResolverList []string `json:"ResolverList"`
	//Opt   Options
}

type Options struct {
	IPv4only bool `json:"IPv4only"`
	IPv6only bool `json:"IPv6only"`
}

func Init(log logger.Logger, zc Map[Zone], sc Map[Server]) Config {
	var conf Config
	conf.Log = log
	conf.Zones = zc
	conf.Cache = sc
	conf.IPv4only = true

	var root Zone
	root.Preload("root-hints.json")
	conf.Zones.Set(".", root)

	conf.ResolverList = []string{"1.1.1.1", "8.8.8.8", "8.8.4.4", "9.9.9.9"}

	return conf

}

func (c *Config) GetResolver() string {
	return c.ResolverList[rand.Intn(len(c.ResolverList))]
}

// if full is true every NS in delegation will be asked for zone
// else, will only collect data from first available server
func PrepZone(name string, cfg *Config) (Zone, error) {

	var zone Zone
	cfg.Log.Debug("Prepping zone", "zone", name)
	// Try to get zone from concurrent map
	if zone, ok := cfg.Zones.Get(name); ok {
		cfg.Log.Debug("Found zone in cache", "zone", name)
		// If the zone is fully primed (200), return it.
		if zone.Status == 200 {
			cfg.Log.Debug("Zone ready", "zone", name, "status", strconv.FormatInt(int64(zone.Status), 10))
			return zone, nil
		}

		// Otherwise, start checking and adding info to zone object
		cfg.Log.Debug("Zone not ready", "zone", name, "status", strconv.FormatInt(int64(zone.Status), 10))
	}
	// If zone is not in cache at all, create a new zone and try to populate it
	cfg.Log.Debug("Creating placeholder for zone", "zone", name)
	zone.Name = name
	zone.Status = 201

	// Get a list of the parent zone nameservers to query for delegation data
	// remove leftmost label to get name of parent zone
	parentZoneName := StripLabelFromLeft(zone.Name)
	nslist, err := Nameservers(parentZoneName, cfg)

	if err != nil {
		return zone, err
	}

	// Populate the parent nameserver info
	err = zone.QueryParentForDelegation(nslist, cfg)
	if err != nil {
		cfg.Log.Debug("Error doing QueryParentForDelegation()", "ERROR", err)
	}

	// Check if the Zone exists, is an actual zone or if it is a hostname/empty non-terminal

	// Populate the zones nameserver info
	err = zone.QuerySelfForNS(cfg)
	if err != nil {
		cfg.Log.Debug("Error doing QuerySelfForNS()", "ERROR", err)
	}

	//zone.Status = 200

	return zone, err

}

// Move to Zones
func Nameservers(ZoneName string, cfg *Config) (map[string]string, error) {
	// Try to get the parent zone from cache
	cfg.Log.Debug("Loading zone", "zone", ZoneName)

	if pzone, ok := cfg.Zones.Get(ZoneName); ok {
		cfg.Log.Debug("Parent zone found", "zone", ZoneName)

		switch pzone.Status {
		case 200:
			// All is going smoothly
			cfg.Log.Debug("Zone ready", "zone", pzone.Name, "status", strconv.FormatInt(int64(pzone.Status), 10))
		case 204:
			// Not a proper zone. Find out the TrueParentZone
			cfg.Log.Debug("Not a proper zone", "zone", pzone.Name, "status", strconv.FormatInt(int64(pzone.Status), 10))
			cfg.Log.Debug("Trying to get data from True Parent Zone", "TPZ", pzone.InZone)
			if tpz, ok := cfg.Zones.Get(pzone.InZone); ok {
				pzone = tpz
			}
		default:
			cfg.Log.Debug("Fell through to default in switch @ func Namesevsers()", "status", strconv.FormatInt(int64(pzone.Status), 10))
		}

		// Return set of NS, if in there are any
		if cfg.IPv4only {
			return pzone.GetNSIP4(), nil
			cfg.Log.Debug("LIST", "IPv4", pzone.GetNSIP4())
		}
		if cfg.IPv4only {
			return pzone.GetNSIP6(), nil
		}

		return pzone.GetNSIP(), nil

	}

	return map[string]string{}, fmt.Errorf("Cache not properly prepared")

}
