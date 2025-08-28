package dig

import (
	"strings"
	"zonetree/logger"

	"github.com/miekg/dns"
)

//var log = logger.PrintDebugLog()

func QueryForDelegation(q Query, log logger.Logger) {

	msg, err := Dig(q)
	if err != nil {
		log.Debug("Error looking up domain", "domain", err.Error())
	}

	rcode := dns.RcodeToString[msg.MsgHdr.Rcode]

	if rcode == "NOERROR" {
		log.Debug("Got reply", "QNAME", q.Qname, "server", q.Nameserver)
	}
}

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

func QndQuery(qname, resolver string, log logger.Logger) ([]string, error) {

	var iplist []string

	q := NewQuery()
	q.RD = true
	q.Nameserver = resolver
	q.Qname = qname

	// Get IPv4 servers
	q.Qtype = "A"

	log.Debug("Sending query", "Query", q)

	msg, err := Dig(q)

	if err != nil {
		log.Error("Error doing QndQuery (A) ", "domain", err.Error())
	}

	rcode := dns.RcodeToString[msg.MsgHdr.Rcode]

	if rcode == "NOERROR" {
		for _, an := range msg.Answer {
			iplist = append(iplist, dns.Field(an, 1))
		}

	}

	// Get IPv6 servers
	q.Qtype = "AAAA"
	msg, err = Dig(q)

	if err != nil {
		log.Error("Error doing QndQuery (AAAA)", "domain", err.Error())
	}

	rcode = dns.RcodeToString[msg.MsgHdr.Rcode]

	if rcode == "NOERROR" {
		for _, an := range msg.Answer {
			iplist = append(iplist, dns.Field(an, 1))
		}

	}

	return iplist, err
}
