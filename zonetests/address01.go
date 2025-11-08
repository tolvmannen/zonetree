package zonetests

import (
	//"fmt"
	//"net/netip"
	"zonetree/cache"
)

// Address01
//
// Struct that implements the ZoneTest interface
type Address01 struct {
	zone string       // name of zone to test
	UNS  []cache.NSIP // nameserver info when running an undelegated test
	//DS       string       // DS record when running undelegated test
	status   int8 // See TestStatusMap for ref.
	tries    int8
	messages []string // Key for message table
	cfg      *cache.Config
}

func (a *Address01) Status() int8 { return a.status }

func (a *Address01) New(cfg *cache.Config, zi ZoneInfo) error {
	a.zone = zi.Name
	a.cfg = cfg

	return nil
}

/*

func (a *Address01) Run() {

	if len(a.UNS) > 0 {

	}


	for _,block := range a.cfg.SuIP {
		return
	}

}

*/
