package dig

import (
	"github.com/miekg/dns"
	"net"
	"os"
	"strings"
)

// StripLabelFromLeft
//
// Returns the parent domain/zone of the given name
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

// ToFQDN
//
// Wrapper for dns.Fqdn
func ToFQDN(name string) string {
	return dns.Fqdn(name)
}

// Path
//
// Returns a lookup path from root to leaf of the domain tree.
// ["se","examples.se","www.examples.se"]
func Path(dom string) []string {
	var tree []string
	dom = dns.Fqdn(dom)
	labels := strings.Split(dom, ".")
	for k := range labels {
		l := strings.Join(labels[k:], ".")
		if l != "" {
			tree = append([]string{l}, tree...)
		}
	}
	return tree
}

// GetSystemResolver
//
// If there is a need to fall back on the system resolver
func GetSystemResolver(ipver string) string {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		//fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	// check for the first available server of right i version
	var ns string
	for _, ip := range conf.Servers {
		ns = net.ParseIP(ip).String()
		if ipver == "6" && strings.Contains(ns, ":") {
			break
		}
		if ipver == "4" && strings.Contains(ns, ".") {
			break
		}

	}
	return ns

}
