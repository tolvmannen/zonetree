package dig

import (
	"github.com/miekg/dns"
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
