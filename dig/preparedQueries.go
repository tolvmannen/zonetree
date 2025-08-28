package dig

// var DefaultDNSquery = Query{
func NewQuery() Query {
	return Query{
		Nameserver: "",
		Transport:  "udp",
		Qname:      "",
		Qtype:      "",
		Port:       "53",
		IpVersion:  "4",
		AA:         false,
		AD:         false,
		CD:         false,
		RD:         false,
		DO:         false,
		NoCrypto:   false,
		Nsid:       false,
		ShowQuery:  false,
		UDPsize:    1232,
		Tsig:       "",
	}
}

var GetPatentData = Query{
	Nameserver: "",
	Transport:  "udp",
	Qname:      "",
	Qtype:      "A",
	Port:       "53",
	IpVersion:  "4",
	AA:         false,
	AD:         false,
	CD:         false,
	RD:         false,
	DO:         true,
	NoCrypto:   false,
	Nsid:       false,
	ShowQuery:  false,
	UDPsize:    1232,
	Tsig:       "",
}

func SoaQuery() Query {
	return Query{
		Nameserver: "",
		Transport:  "udp",
		Qname:      "",
		Qtype:      "SOA",
		Port:       "53",
		IpVersion:  "4",
		AA:         false,
		AD:         false,
		CD:         false,
		RD:         false,
		DO:         true,
		NoCrypto:   false,
		Nsid:       false,
		ShowQuery:  false,
		UDPsize:    1232,
		Tsig:       "",
	}
}

func ParentDelegQuery(parent, child string) Query {
	return Query{
		Nameserver: parent,
		Transport:  "udp",
		Qname:      child,
		Qtype:      "NS",
		Port:       "53",
		IpVersion:  "4",
		AA:         false,
		AD:         false,
		CD:         false,
		RD:         false,
		DO:         true,
		NoCrypto:   false,
		Nsid:       false,
		ShowQuery:  false,
		UDPsize:    1232,
		Tsig:       "",
	}
}
