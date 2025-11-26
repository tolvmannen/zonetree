package dig

// TODO Rename to query?

import (
	"fmt"
	"net"
	"strings"

	"zonetree/logger"
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
	SendQuery  func(Query, logger.Logger) (DigData, error)
}

// Sanitize
//
// sanitize input data as precaution
func (q *Query) Sanitize() {
	q.Transport = strings.ToLower(q.Transport) // needs to be lower case.
}

// Send
//
// If no SendQuery function is set this will default to using
// the general SendQuery (dig.Dig) function
func (q *Query) Send(log logger.Logger) (DigData, error) {
	if q.SendQuery == nil {
		return SendQuery(*q, log)
	}
	return q.SendQuery(*q, log)

}

// GetLookupNS
//
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
