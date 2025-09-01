package dig

import (
	"strings"
	"zonetree/logger"

	"github.com/miekg/dns"
)

// Setting up some custom structs to more easily
// handle the data retured from a query
// Also, removes the dependency on miekg/dns in
// other parts
type DigData struct {
	Rcode         string
	Answer        []DigRR
	Authoritative []DigRR
	Additional    []DigRR
}

type DigRR struct {
	Name  string
	Rtype string
	Ttl   uint32
	Rdata []string
}

// Returnd Rdata as one (1) space separated string
func (rr *DigRR) GetRdata() string {
	return strings.Join(rr.Rdata, " ")
}

// Return Rdata as separate fields
func (rr *DigRR) GetRdataFields() []string {
	return rr.Rdata
}

func GetDelegation(q Query, log logger.Logger) (DigData, error) {

	var data DigData

	msg, err := Dig(q)
	if err != nil {
		log.Error("Nameserver reported error looking up domain", "domain", err.Error())
		return data, err
	}

	data.Rcode = dns.RcodeToString[msg.MsgHdr.Rcode]

	if data.Rcode == "NOERROR" {
		log.Debug("Got reply", "QNAME", q.Qname, "server", q.Nameserver)

		// Go through all the sections of the response and
		// sort the right info into the DigData struct
		for _, au := range msg.Answer {
			var rr DigRR
			head := *au.Header()
			rr.Rtype = dns.Type(head.Rrtype).String()
			rr.Name = head.Name
			rr.Ttl = head.Ttl
			for i := 1; i <= dns.NumField(au); i++ {
				rr.Rdata = append(rr.Rdata, dns.Field(au, i))
			}
			data.Answer = append(data.Answer, rr)
		}
		for _, au := range msg.Ns {
			var rr DigRR
			head := *au.Header()
			rr.Rtype = dns.Type(head.Rrtype).String()
			rr.Name = head.Name
			rr.Ttl = head.Ttl
			for i := 1; i <= dns.NumField(au); i++ {
				rr.Rdata = append(rr.Rdata, dns.Field(au, i))
			}
			data.Authoritative = append(data.Authoritative, rr)
		}
		for _, au := range msg.Extra {
			var rr DigRR
			head := *au.Header()
			rr.Rtype = dns.Type(head.Rrtype).String()
			rr.Name = head.Name
			rr.Ttl = head.Ttl
			for i := 1; i <= dns.NumField(au); i++ {
				rr.Rdata = append(rr.Rdata, dns.Field(au, i))
			}
			data.Additional = append(data.Additional, rr)
		}

	}

	return data, err
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
