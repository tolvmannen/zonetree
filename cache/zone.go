package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"zonetree/dig"

	"github.com/miekg/dns"
)

// Having a laugh with HTTP status code references and stuff...
var ZoneStatus = map[int32]string{
	69:  "Nice", // Just nice
	200: "Zone OK",
	201: "Zone placeholder created",
	204: "Not a Zone",      // Either a hostname or an empty non-terminal
	206: "Zone incomplete", // Zone have not yet been, or could not be, entierly processed
	207: "Zone OK-ish?",    // Multi-Status - No consensus on status
	403: "REFUSED",
	404: "NXDOMAIN",
	420: "I'm hallucinating", // What even is this?
	422: "Unable to comply",  // Unprocessable. Most likely disallowed due to Config option.
	500: "Server error",
}

// Zone
//
// This struct holds all relevant data for a zone.
type Zone struct {
	Name     string     `json:"Name"`     // Name of the Zone
	ZoneNS   []ZoneNS   `json:"NS"`       // All NS from all instances of the Authoritative name servers
	InZone   string     `json:"InZone"`   // Zone the record belongs to (in case of hosts and empty non-terminals)
	ParentNS []ParentNS `json:"ParentNS"` // All NS in all instances from the name servers of the Parent Zone (i.e. delegations)
	NSIP     []NSIP     `json:"NSIP"`     // All NS Name <-> IP pairs found in both delegation and in Authoritative name servers
	Status   int32      `json:"Status"`   // See ZoneStatus
}

// ZoneNS
//
// Struct to hold relevant data for the Zones Authoritative nameservers
type ZoneNS struct {
	Self   int8     `json:"Self"` // Reference to index in the Zone struct NSIP list from where the data was received
	NS     []int8   `json:"NS"`   // Reference to indexes in the Zone struct NSIP list containing NS record info
	SOA    string   `json:"SOA"`
	DNSKEY []string `json:"DNSKEY"`
	RRSIG  []string `json:"RRSIG"`
}

// ParentNS
//
// Struct to hold relevant data from name servers of Parent zone
type ParentNS struct {
	Name        string   `json:"Name"`
	IP          string   `json:"IP"`
	NS          []int8   `json:"NS"` // Reference to indexes in the Zone struct NSIP list containing NS record info
	DS          []string `json:"DS"`
	RRSIG       []string `json:"RRSIG"`
	ChildStatus int32    `json:"ChildStatus"` // Used to keep track of inconsitencies in delgation NS set @ parents
}

// NSIP
//
// Struct to hold relevant NS / Delegation data
type NSIP struct {
	Name       string `json:"Name"`
	IP         string `json:"IP"`
	ZoneStatus int32  `json:"ZoneStatus"` // Status of the zone according to this server
}

// Server
//
// Struct used for keeping relevant information on nameservers (Resolvers and Authoritative)
// in a global cache
type Server struct {
	IP []string `json:"IP"`
}

// Get the self referenced entries from the NSIP list for the Zone
func (z *Zone) GetNSIP() map[string]string {
	var nsset = make(map[string]string)
	for _, i := range z.ZoneNS {
		nsset[z.NSIP[i.Self].IP] = z.NSIP[i.Self].Name
	}
	return nsset
}

func (z *Zone) GetNSIP4() map[string]string {
	var nsset = make(map[string]string)
	for _, i := range z.ZoneNS {
		if strings.Count(z.NSIP[i.Self].IP, ":") < 1 {
			nsset[z.NSIP[i.Self].IP] = z.NSIP[i.Self].Name
		}
	}
	return nsset
}

func (z *Zone) GetNSIP6() map[string]string {
	var nsset = make(map[string]string)
	for _, i := range z.ZoneNS {
		if strings.Count(z.NSIP[i.Self].IP, ":") > 0 {
			nsset[z.NSIP[i.Self].IP] = z.NSIP[i.Self].Name
		}
	}
	return nsset
}

/*
// Check if all parent nameservers are in agreement of the status of the
// child zone
func (z *Zone) ChildStatusConsensus() bool {
	cs := make(map[int]int) // map[status]counter
	for _, p := range z.ParentNS {
		cs[int(p.ChildStatus)]++
	}
	// if there is exactly one status, all servers are in agreement
	// return true and the agreed upon status
	if len(cs) == 1 {
		return true
	}
	return false
}
*/

func BuildZoneCache(z string, cfg *Config) {

	// If asked to check . (i.e. ROOT zone)
	// do nothing 4 now
	if z == "." {
		return
	}
	// Make DNS tree list
	list := dig.Path(z)

	// Loop through the nodes and perp the zone.
	// Will start at TLD, because ROOT should already be primed,
	// or the program would have blown up by now.
	for _, node := range list {
		// If not specifically checking a tld, don't query ALL the root servers
		// for TLD delegation
		zone, err := PrepZone(node, cfg)
		//zone, err := MockZone(node, full)

		if err != nil {
			cfg.Log.Error("Error preparing zone", "zone", node, "Error", err)
		}
		cfg.Zones.Set(node, zone)
	}

}

// func (z *Zone) DelegationInBailiwick(ns string) bool {
func DelegationInBailiwick(nsname, dom string) bool {
	// make fqdn and compare the domain w the last part of NS name.
	//dom := dns.Fqdn(ns[len(ns)-len(z.Name):])
	domain := dns.Fqdn(nsname[len(nsname)-len(dom):])
	if domain == nsname {
		return true
	}
	return false
}

func (z Zone) Print() {
	jz, err := json.Marshal(z)
	if err != nil {
		fmt.Printf("Warning: unable to marshal(%s)\n", err.Error())
	}
	fmt.Printf("\n%+v\n", string(jz))
}

func (z Zone) ToJson() (string, error) {
	jz, err := json.Marshal(z)
	return string(jz), err
}

func (z Zone) ToPrettyJson() (string, error) {
	jz, err := json.MarshalIndent(z, "", "  ")
	return string(jz), err
}

// preload zone data from json file
// mainly for root-hints and se-hints
func (z *Zone) Preload(file string) {

	js, err := os.ReadFile("hints/" + file)
	if err != nil {
		fmt.Printf("Warning: unable to load hint file %s (%s)\n", file, err.Error())
	}
	json.Unmarshal(js, z)
	if err != nil {
		fmt.Printf("Warning: unable to unmarshal(%s)\n", err.Error())
	}

}

func NewZoneCache() Map[Zone] {
	return NewMapFromConfig[Zone](false) // set true for dummy
}

func NewServerCache() Map[Server] {
	return NewMapFromConfig[Server](false) // set true for dummy
}

// return NS data for all nameservers in a namserver delegation.
func (z *Zone) QueryParentForDelegation(nslist map[string]string, cfg *Config) error {

	q := dig.NewQuery()
	q.Qname = z.Name
	// query for SOA and set DO (qmin-ish and may save a query or two)
	q.Qtype = "SOA"
	q.DO = true

	for ip, name := range nslist {

		// Check if the IP is already in the Delegation NS set of the zone
		pid := slices.IndexFunc(z.ParentNS, func(ns ParentNS) bool {
			return ns.IP == ip
		})

		// If not, create ann entry ad add it
		if pid < 0 {
			var parent ParentNS
			parent.IP = ip
			parent.Name = name
			z.ParentNS = append(z.ParentNS, parent)
			// Then fetch id for later operations
			pid = slices.IndexFunc(z.ParentNS, func(ns ParentNS) bool {
				return ns.IP == ip
			})
			cfg.Log.Debug("DELEGATION: IP not found in ParentNS. Creating placeholder", "IP", ip, "Name", name, "ID", pid)
		} else {
			cfg.Log.Debug("DELEGATION: IP already in ParentNS.", "IP", ip, "ID", pid)
		}

		q.Nameserver = ip
		cfg.Log.Debug("Parent Query:", "query", q)
		msg, err := dig.GetDelegation(q, cfg.Log)
		if err != nil {
			//cfg.Log.Error("DELEGATION: Error looking up domain", "domain", err.Error())
			cfg.Log.Debug("Query Failed Trying next nameserver in list", "ERROR", err)
			continue
		}

		if msg.Rcode == "NOERROR" {

			cfg.Log.Debug("DELEGATION: NOERROR", "QNAME", q.Qname, "server", q.Nameserver)

			// Get info from Auth section
			// Extract DNSSEC info, if any, and make a list of delegation
			// Name Servers
			var delegns []NSIP
			for _, au := range msg.Authoritative {

				// RDATA is in dns.RR.<section>[1:]
				switch au.Rtype {
				case "NS":
					// create placeholder NS struct to put IP in later
					name := au.GetRdata()
					// Check if the name is already in the NSIP list of the zone
					id := slices.IndexFunc(delegns, func(ns NSIP) bool {
						return ns.Name == name
					})
					// If not, create ann entry and add it
					if id < 0 {
						var nsip NSIP
						nsip.Name = name
						delegns = append(delegns, nsip)
					}
					//z.ParentNS[pid].ChildStatus = 200
				case "DS":
					z.ParentNS[pid].DS = append(z.ParentNS[pid].DS, au.GetRdata())
				case "RRSIG":
					z.ParentNS[pid].RRSIG = append(z.ParentNS[pid].RRSIG, au.GetRdata())
				case "SOA":
					// NORROR + Authoritative answer + SOA in Authoritative section
					// indicates that name in either a host name or an empty non-terminal
					// Set statuses accordingly and make a note of true parent zone
					z.ParentNS[pid].ChildStatus = 204
					z.Status = 204
					z.InZone = au.Name
					cfg.Log.Debug("[Parent] reported [Name] to be a part of [Zone]", "Parent", q.Nameserver, "Name", q.Qname, "Zone", au.Name)

				default:
				}

			}

			// Get all glue that is provided, but dont trust it to be complete.
			// This will save a few lookups further on
			for _, e := range msg.Additional {
				if e.Rtype == "A" || e.Rtype == "AAAA" {
					// check if an identical entry exists
					id := slices.IndexFunc(delegns, func(ns NSIP) bool {
						return ns.IP == e.GetRdata() && ns.Name == e.Name
					})
					if id < 0 {
						// check for entry with name but no IP
						id := slices.IndexFunc(delegns, func(ns NSIP) bool {
							return ns.Name == e.Name && ns.IP == ""
						})

						if id < 0 {
							var nsip NSIP
							nsip.IP = e.GetRdata()
							nsip.Name = e.Name
							delegns = append(delegns, nsip)
						} else {
							delegns[id].IP = e.GetRdata()
						}
					}
				}
			}

			cfg.Log.Debug("DELEGATION: Pepared ns <-> ip list for Parent NS", "NS", q.Nameserver, "LIST", delegns)
			// Go through the list of ns <-> ip and check if there is
			// already an identical entry in the zones NSIP list.
			// If so, add a reference. If not, add both entry and reference
			for _, e := range delegns {
				// No IP here means it was not in Glue.
				if e.IP != "" {
					id := slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
						return ns.Name == e.Name && ns.IP == e.IP
					})
					if id < 0 {

						var nsip NSIP
						nsip.IP = e.IP
						nsip.Name = e.Name
						z.NSIP = append(z.NSIP, nsip)
						// Since we don't know the index it got when
						//inserted, we need to fetch it
						id = slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
							return ns.Name == e.Name && ns.IP == e.IP
						})
						cfg.Log.Debug("DELEGATIO|: Did not find identical match. Adding new entry", "NSIP", nsip, "Index", id)
					} else {
						cfg.Log.Debug("DELEGATION: Found identical match.", "NS", e.Name, "IP", e.IP, "Index", id)
					}
					z.ParentNS[pid].NS = append(z.ParentNS[pid].NS, int8(id))
				} else {
					cfg.Log.Debug("DELEGATION: No IP in ns <-> pair. Doing recursive lookup", "Name", e.Name)

					var iplist []string
					// Check if the name server is in the global cache
					if server, ok := cfg.Cache.Get(e.Name); ok {
						cfg.Log.Debug("DELEGATION: Nameserver found in global cache", "Name", e.Name)
						for _, ip := range server.IP {
							iplist = append(iplist, ip)
						}
					}

					if len(iplist) < 1 {
						cfg.Log.Debug("DELEGATION: Nameserver NOT in global cache. Querying resolver.", "Name", e.Name)
						// Cheat and use a resolver to get the IP(s) for the NS name
						iplist, _ = dig.QndQuery(e.Name, cfg.GetResolver(), cfg.Log)
						if len(iplist) > 0 {
							server := Server{IP: iplist}
							cfg.Cache.Set(e.Name, server)
						}
					}

					for _, ip := range iplist {
						// Even if the IP was not in the Glue for this NS
						// it might have been added when processing another
						// nameserver. Extra check just in case.
						cfg.Log.Debug("DELEGATION: IP-LIST for nameserver.", "Name", e.Name, "IP", ip)
						id := slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
							return ns.Name == e.Name && ns.IP == ip
						})
						if id < 0 {
							cfg.Log.Debug("", "ID", id, "name", e.Name, "ip", ip)
							var nsip NSIP
							nsip.IP = ip
							nsip.Name = e.Name
							z.NSIP = append(z.NSIP, nsip)
							id = slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
								return ns.Name == e.Name && ns.IP == ip
							})
							cfg.Log.Debug("DELEGATION: ns <-> ip pair (still) not in list. Adding new entry", "NSIP", nsip, "Index", id)
						}
						z.ParentNS[pid].NS = append(z.ParentNS[pid].NS, int8(id))
					}
				}
			}

			// If not already set, set Child Zone status at parent level to OK
			if z.ParentNS[pid].ChildStatus == 0 {
				z.ParentNS[pid].ChildStatus = 200
			}
			// Sort the NS entries for easier comparison later
			slices.Sort(z.ParentNS[pid].NS)

		}

		if msg.Rcode == "NXDOMAIN" {
			// If the zone can't be found at the parent NS
			// set status accordingly
			z.ParentNS[pid].ChildStatus = 404
		}

		if msg.Rcode == "REFUSED" {
			// If the the parent NS refuses to answer
			// set status accordingly
			z.ParentNS[pid].ChildStatus = 403
		}
	}

	// If the zone was created (201), update status reflect that ParenNS is populated
	// If another status has been set, keep that.
	if z.Status == 201 {
		z.Status = 206 // Partially Primed zone
	}

	return nil
}

// Query all nameservers that we've found in delegations form paretn nameservers.
func (z *Zone) QuerySelfForNS(cfg *Config) error {

	q := dig.NewQuery()
	q.Qname = z.Name
	// query for SOA and set DO (qmin-ish and may save a query or two)
	q.Qtype = "NS"
	q.DO = true

	// full set should be in z.NSID
	for i, nsip := range z.NSIP {

		cfg.Log.Debug("Querying server", "nr", i+1, "of", len(z.NSIP), "in list", nsip)

		// Dont query IP-addresses of the wrong version if the option to
		// use only 4 or 6 is set.
		if cfg.IPv4only && strings.Contains(nsip.IP, ":") {
			cfg.Log.Debug("IPv4 only. Ignoring address).", "IP", nsip.IP)
			z.NSIP[i].ZoneStatus = 422 // won't do the v6 for conf reasons
			continue
		}

		if cfg.IPv6only && !strings.Contains(nsip.IP, ":") {
			cfg.Log.Debug("IPv6 only. Ignoring address).", "IP", nsip.IP)
			z.NSIP[i].ZoneStatus = 422 // won't do the v4 for conf reasons
			continue
		}

		q.Nameserver = nsip.IP

		cfg.Log.Debug("SELF Query:", "query", q)
		msg, err := dig.GetDelegation(q, cfg.Log)
		if err != nil {
			z.NSIP[i].ZoneStatus = 500
			cfg.Log.Debug("Query Failed Trying next nameserver in list", "ERROR", err)
			continue
		}

		/*
			msg, err := dig.Dig(q)
			if err != nil {
				cfg.Log.Error("Error looking up domain", "domain", err.Error())
			}

			rcode := dns.RcodeToString[msg.MsgHdr.Rcode]
		*/

		rcode := msg.Rcode

		if rcode == "NOERROR" {

			// So far zone is OK.
			z.NSIP[i].ZoneStatus = 200

			// We're expecting an Authoritative answer
			// If not AA go tonext server
			if !msg.AA {
				cfg.Log.Debug("Got NON-AUTHORITATIVE reply. Proceeding to next server", "QNAME", q.Qname, "server", q.Nameserver)
				continue
			}

			cfg.Log.Debug("Got reply", "QNAME", q.Qname, "server", q.Nameserver)

			if len(msg.Answer) < 1 {
				cfg.Log.Debug("Answer section empty")
				// Only log this 4 now
			}

			if len(msg.Authoritative) < 1 {
				cfg.Log.Debug("Authoritative section empty")
				// Only log this 4 now
			}

			if len(msg.Additional) < 1 {
				cfg.Log.Debug("Additional section empty")
				// Only log this 4 now
			}

			var zns ZoneNS

			// nameservers in NS section
			// This will be used to get IP addresses for nameservers
			// Not found in glue and/or not in bailiwick
			var nsrr []string
			for _, an := range msg.Answer {

				// RDATA is in dns.RR.<section>[1:]
				if an.Rtype == "NS" {
					nsrr = append(nsrr, an.GetRdata())
				}
			}

			// Get all glue that is provided, but dont trust it to be complete.
			// Add any missing entries to the NSIP list and remove tne name from
			cfg.Log.Debug(" -- Parsing Additional section --")
			for _, e := range msg.Additional {
				// RDATA is in dns.RR.<section>[1:]
				if e.Rtype == "A" || e.Rtype == "AAAA" {
					cfg.Log.Debug("Looking for  existing entry in Zone NSID list", "Name", e.Name, "IP", e.GetRdata())
					id := slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
						return ns.Name == e.Name && ns.IP == e.GetRdata()
					})

					if id < 0 {
						var nsip NSIP
						nsip.IP = e.GetRdata()
						nsip.Name = e.Name
						z.NSIP = append(z.NSIP, nsip)
						cfg.Log.Debug("No Name in z.NSIP, adding entry", "Name", nsip.Name, "IP", nsip.IP)
						// Get the index of the inserted NSID
						id = slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
							return ns.Name == e.Name && ns.IP == e.GetRdata()
						})
					}
					// Add the id as a NSID reference in the ZoneNS.
					cfg.Log.Debug("Adding reference to NS list", "ID", id)
					zns.NS = append(zns.NS, int8(id))
					// Add a self reference if the IP matches that of the queried
					// nameserrver
					if nsip.IP == e.GetRdata() {
						cfg.Log.Debug("Adding SELF reference", "My IP", e.GetRdata(), "Queried IP", nsip.IP)
						zns.Self = int8(id)
					}

					rrid := slices.Index(nsrr, e.Name)
					if rrid > -1 {
						cfg.Log.Debug("Removing name from NSRR list", "Name", e.Name)
						nsrr = slices.Delete(nsrr, rrid, rrid+1)

					}
				}
			}

			cfg.Log.Debug("Finding IP for unresolved NS names", "NSRR", nsrr)
			for _, name := range nsrr {

				// TODO Contemplate order of checking bailiwick, then cache
				// vs the other way around

				var iplist []string

				if DelegationInBailiwick(name, z.Name) {
					cfg.Log.Debug("Making Biliwick Lookup", "Name", name)
					iplist, err = dig.QndQuery(name, nsip.IP, cfg.Log)
				}
				if err != nil {
					cfg.Log.Error("Error in Biliwick Lookup", "ERR", err)
				}
				// If the delegation is out of bailiwick or if something
				// went wrong and the Authoritative NS couldn't
				// provide a lookup, look in cache for server.
				//

				cfg.Log.Debug("IP-list before cache", "list", iplist)
				if len(iplist) < 1 {
					cfg.Log.Debug("Making Cache Lookup", "Name", name)
					if server, ok := cfg.Cache.Get(name); ok {
						iplist = append(iplist, server.IP...)
					}
				}
				cfg.Log.Debug("IP-list after cache", "list", iplist)

				// If that fails, use a resolver to get the IP(s) for
				// the NS name
				if len(iplist) < 1 {
					cfg.Log.Debug("Making Resolver Lookup", "Name", name)
					iplist, _ = dig.QndQuery(name, cfg.GetResolver(), cfg.Log)
					// if this succeeds, save server in global cache
					if len(iplist) > 0 {
						server := Server{IP: iplist}
						cfg.Cache.Set(name, server)
					}
				}
				for _, ip := range iplist {
					// Even if the IP was not in the Glue for this NS
					// it might have been added when processing another
					// nameserver. Extra check just in case.
					cfg.Log.Debug("IP-LIST for nameserver.", "Name", name, "IP", ip)
					id := slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
						return ns.Name == name && ns.IP == ip
					})
					cfg.Log.Debug("Re-checking zones NSID list", "ID", id, "name", name, "ip", ip)
					if id < 0 {
						var nsip NSIP
						nsip.IP = ip
						nsip.Name = name
						z.NSIP = append(z.NSIP, nsip)
						id = slices.IndexFunc(z.NSIP, func(ns NSIP) bool {
							return ns.Name == name && ns.IP == ip
						})
						cfg.Log.Debug("ns <-> ip pair (still) not in list. Adding new entry", "NSIP", nsip, "Index", id)
					}

					// Add the id as a NSID reference in the ZoneNS.
					cfg.Log.Debug("Adding reference to NS list", "ID", id)
					zns.NS = append(zns.NS, int8(id))
					// Add a self reference if the IP matches that of the queried
					// nameserrver
					if nsip.IP == ip {
						cfg.Log.Debug("Adding SELF reference", "My IP", ip, "Queried IP", nsip.IP)
						zns.Self = int8(id)
					}

				}

			}

			//If there is at least 1 working nameserver, aet zone tatus
			// to OK
			z.Status = 200

			// Sort the NS list for easier comparison later
			slices.Sort(zns.NS)
			// Add the ZonNS to the Zone
			z.ZoneNS = append(z.ZoneNS, zns)

		}

		if rcode == "NXDOMAIN" {
			z.NSIP[i].ZoneStatus = 404
		}

		if rcode == "REFUSED" {
			z.NSIP[i].ZoneStatus = 403
		}

	}

	return nil
}

func StripLabelFromLeft(z string) string {

	var parent string

	// Split domain int labels and move up one step
	labels := dns.SplitDomainName(z)
	// ROOT is parent to itself and TLDs
	if len(labels) < 2 {
		parent = "."
	} else {
		parent = dns.Fqdn(strings.Join(labels[1:], "."))
	}
	return parent
}

func MakeFQDN(name string) string {
	return dns.Fqdn(name)
}
