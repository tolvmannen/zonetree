package dig

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Query struct {
	Nameserver string `json:"Nameserver"`
	Transport  string `json:"Transport"`
	Qname      string `json:"Qname"`
	Qtype      string `json:"Qtype"`
	Port       string `json:"Port"`
	IpVersion  string `json:"IpVersion"`
	AA         bool   `json:"AA"`
	AD         bool   `json:"AD"`
	CD         bool   `json:"CD"`
	RD         bool   `json:"RD"`
	DO         bool   `json:"DO"`
	NoCrypto   bool   `json:"NoCrypto"`
	Nsid       bool   `json:"Nsid"`
	ShowQuery  bool   `json:"ShowQuery"`
	UDPsize    uint16 `json:"UDPsize"`
	Tsig       string `json:"Tsig"`
}

type DigOut struct {
	Qname      string        `json:"Qname"`
	Query      *dns.Msg      `json:"Query"`    // message sent to name server
	Response   *dns.Msg      `json:"Response"` // response from nameserver
	RTT        time.Duration `json:"Round trip time"`
	Nameserver string        `json:"Nameserver"` // Name server IP
	QNSname    string        `json:"QNSname"`    // resolver name before translation
	ShowQuery  bool          `json:"ShowQuery"`
	MsgSize    int           `json:"Message Size"`
	Transport  string        `json:"Transport"`
}

// sanitize input data as precaution
func (q *Query) Sanitize() {
	q.Transport = strings.ToLower(q.Transport) // needs to be lower case.
}

/*
func GetSystemResolver() string {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	ns := conf.Servers[0]
	// Strip the [ and ] from around the nameserver obtained from /etc/resolv.conf
	//ns = ns[1 : len(ns)-1]
	return ns
}
*/

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

// Harmonize lookup nameserver to always use IP:Port
// Check if valid IP. If not assume, hostname and look it up, selecting the first available ip
// of correct version
func (q *Query) GetLookupNS() string {
	var ns string

	// If no nameserver was passed, use system resolver
	if len(q.Nameserver) == 0 {
		ns = GetSystemResolver("4")
		return ns
	}

	ip := net.ParseIP(q.Nameserver)
	if ip != nil {
		if q.IpVersion == "6" {
			ns = "[" + q.Nameserver + "]:" + q.Port
		} else {
			ns = q.Nameserver + ":" + q.Port
		}
	} else {
		IPlist, err := net.LookupIP(q.Nameserver)
		if err != nil {
			fmt.Printf("Nameserver lookup error: %v\n", err)
		} else {
			for _, ip := range IPlist {
				if q.IpVersion == "6" {
					if strings.Count(ip.String(), ":") >= 2 {
						ns = "[" + ip.String() + "]:" + q.Port
						break
					}
				} else {
					// If address contains more than 1 ':', it's a V6 address. Go next.
					if strings.Count(ip.String(), ":") < 2 {
						ns = ip.String() + ":" + q.Port
						break
					}
				}
			}
		}

	}
	return ns
}
