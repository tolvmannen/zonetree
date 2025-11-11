package zonetests

import (
	"zonetree/cache"
)

type Basic01 struct {
	zone string       // name of zone to test
	UNS  []cache.NSIP // nameserver info when running an undelegated test
	//DS       string       // DS record when running undelegated test
	status   int8 // See TestStatusMap for ref.
	tries    int8
	messages []string // Key for message table
	cfg      *cache.Config
}

func (t *Basic01) Status() int8 { return a.status }

func (t *Basic01) New(cfg *cache.Config, zi ZoneInfo) error {
	t.zone = zi.Name
	t.cfg = cfg

	return nil
}
