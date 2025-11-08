package cache

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"math/rand/v2"
	"os"
	"strconv"
	"zonetree/logger"
)

type Config struct {
	Log          logger.Logger
	Zones        Map[Zone]
	Cache        Map[Server]
	SuIP         []SuIP   // Special Use IP-adresses
	IPv4only     bool     `json:"IPv4only"`
	IPv6only     bool     `json:"IPv6only"`
	ResolverList []string `json:"ResolverList"`
	Opt          Options  `json:"Opt"`
}

// Options
//
// There are a number of ways of doing lookups using Query Minimization
// Some nameservers use different sequences for adding labels
// All the Qmin-options deal with the Qmin lookup behaviour
// QminLabelSequence	- The number of labels to add, in what sequence.
//
//	Usually the first two (tld, sld) eare added 1 by 1.
//	To save on DNS queries, labels may be added in multiples
//	I.e. 1,1,1,2,2,2,3,3,3... 1,1,1,1,1,1,1... 1,1,3,3,3,3,3... etc
//
// QminSubtractCache	- Count down on cache hits.
// QminStrict		- If true, abort on fail rather than falling back to using the full domain name.
// QminFirstPath	- If true, continue to next label after first successful lookup.
type Options struct {
	IPv4only          bool     `json:"IPv4only" yaml:"IPv4only"`
	IPv6only          bool     `json:"IPv6only" yaml:"IPv6only"`
	ResolverList      []string `json:"ResolverList" yaml:"ResolverList"`
	QminLabelSequence []int8   `json:"QminLabelSequence" yaml:"QminLabelSequence"`
	QminSubtractCache bool     `json:"QminSubtractCache" yaml:"QminSubtractCache"`
	QminStrict        bool     `json:"QminStrict" yaml:"QminStrict"`
	QminFirstPath     bool     `json:"QminFirstPath" yaml:"QminFirstPath"`
}

// Init
//
// Initialses the confguration used when running tests.
func Init() Config {
	var conf Config
	conf.Zones = NewZoneCache()
	conf.Cache = NewServerCache()
	conf.Log = logger.PrintDebugLog()

	conf.SuIP = DefaultSuIP()

	var root Zone
	root.Preload("root-hints.json")
	conf.Zones.Set(".", root)

	conf.DefaultOptions()

	//conf.Opt.QminFirstPath = true

	return conf

}

// DefaultOptions
//
// Default options to use if config file can't be be fond/read
func (c *Config) DefaultOptions() {
	c.Opt = Options{
		IPv4only:          true,
		IPv6only:          false,
		QminSubtractCache: true,
		QminLabelSequence: []int8{1}, // No shortcuts.
		QminStrict:        false,
		QminFirstPath:     false,
		ResolverList:      []string{"1.1.1.1", "8.8.8.8", "8.8.4.4", "9.9.9.9"},
	}

}

// Load
//
// Loads a YAML config file, ovverwriting default options.
func (c *Config) Load(file string) error {

	cf, err := os.ReadFile("profiles/" + file)
	if err != nil {
		fmt.Printf("ReadFile error: %v\n", err)
	}
	yaml.Unmarshal(cf, &c.Opt)
	/*
		if err != nil {
			fmt.Printf("YAML unmarshal error: %v\n", err)
		}
	*/

	return err

}

// RunningConf
//
// Prints out the current loaded conf in YAML format.
func (c *Config) RunningConf() []byte {
	cnf, err := yaml.Marshal(&c.Opt)
	if err != nil {
		fmt.Printf("YAML marshal error: %v\n", err)
	}
	return cnf

}

func (c *Config) GetResolver() string {
	if len(c.Opt.ResolverList) > 0 {
		return c.Opt.ResolverList[rand.IntN(len(c.Opt.ResolverList))]
	}

	// TODO Move this to a default const
	return "1.1.1.1"
}

// PrepZone
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
	nslist, zonecut, err := Nameservers(parentZoneName, cfg)
	zone.ZoneCut = zonecut

	if err != nil {
		return zone, err
	}

	// Populate the parent nameserver info
	// If the option for First Path is set, stop going through the list
	// as soon as enough information to continue down the tree is obtaine
	for ip, name := range nslist {
		pds := zone.QueryParentForDelegation(ip, name, cfg)
		if pds == 200 && cfg.Opt.QminFirstPath {
			break
		}
	}

	status := zone.CalcZoneStatus()

	// If the parent zone has no info about the child zone
	// i.e. 420 it is (most likely) not a proper zone
	// Re-use parents status for the child zone
	if status == 420 {
		if pz, ok := cfg.Zones.Get(StripLabelFromLeft(zone.Name)); ok {
			cfg.Log.Debug("Not proper zone. Re-using status from parent", "Zone", zone.Name, "Parent Zone Status", status)
			status = pz.Status
		}
	}

	zone.Status = status

	// Populate the zones nameserver info
	err = zone.QuerySelfForNS(cfg, cfg.Opt.QminFirstPath)
	if err != nil {
		cfg.Log.Debug("Error doing QuerySelfForNS()", "ERROR", err)
	}

	return zone, err

}

// PrepUndelegatedZone
//
// Undelegated zones relies on user submitted data. No traversal of the
// DNS-tree from ROOT down will be attempted
func PrepUndelegatedZone(name string, cfg *Config, uns []NSIP) (Zone, error) {

	var zone Zone
	var err error
	cfg.Log.Debug("Prepping Undelegated zone from userdata", "zone", name)

	zone.Name = name
	if len(uns) > 0 {
		zone.NSIP = uns
	} else {
		cfg.Log.Debug("Userdata incomplete. Missing NS info", "zone", name)
	}

	return zone, err
}

// TODO Move to Zones
func Nameservers(ZoneName string, cfg *Config) (map[string]string, string, error) {
	// Try to get the zone from cache
	cfg.Log.Debug("Loading zone", "zone", ZoneName)

	// Function returns ZoneCut to make further processing easier.
	var zc string

	if zone, ok := cfg.Zones.Get(ZoneName); ok {
		cfg.Log.Debug("Parent zone found in cache", "zone", ZoneName, "status", zone.Status)

		switch zone.Status {
		case 200:
			// All is going smoothly
			cfg.Log.Debug("Zone ready", "zone", zone.Name, "status", strconv.FormatInt(int64(zone.Status), 10))
		case 204:
			// Not a proper zone. Check ZoneCut
			cfg.Log.Debug("Not a proper zone", "zone", zone.Name, "status", strconv.FormatInt(int64(zone.Status), 10))
			cfg.Log.Debug("Trying to get data from Parent Zone derived from ZoneCut", "ZoneCut", zone.ZoneCut)
			if tpz, ok := cfg.Zones.Get(zone.ZoneCut); ok {
				zone = tpz
			}
		case 404:
			// Not a proper zone. Check ZoneCut
			cfg.Log.Debug("NXDOMAIN", "zone", zone.Name, "status", strconv.FormatInt(int64(zone.Status), 10))
			cfg.Log.Debug("Inherit ZoneCut from closest label", "ZoneCut", zone.ZoneCut)
			if tpz, ok := cfg.Zones.Get(zone.ZoneCut); ok {
				zone = tpz
			}
		default:
			cfg.Log.Debug("Fell through to default in switch @ func Namesevsers()", "status", strconv.FormatInt(int64(zone.Status), 10))
		}

		zc = zone.ZoneCut

		// Return set of NS, if in there are any
		if cfg.Opt.IPv4only {
			return zone.GetNSIP4(), zc, nil
		}
		if cfg.Opt.IPv6only {
			return zone.GetNSIP6(), zc, nil
		}

		return zone.GetNSIP(), zc, nil

	}

	return map[string]string{}, "", fmt.Errorf("Function Nameservers() exploded trying to get nameservers for %s\n", ZoneName)

}
