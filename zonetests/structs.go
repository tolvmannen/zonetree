package zonetests

import (
	"encoding/json"
	"zonetree/cache"
)

/*
type ZoneTests interface {
	Run(param TestParam)
	Passed()   //  boolean based on test criteria
	Messages() // returns list of messags from test results
}
*/

// TestStatusMap
//
// Keeps track on test progression.
var TestStatusMap = map[int8]string{
	0: "Not Started",
	1: "Running",
	2: "Waiting for Retry",
	3: "Disabled",
	4: "Failed",
	5: "Passed",
}

// TestRequest
//
// Client request to run a number of tests on a nmber of zones.
// Delegation information is optional can be supplied on a zone to zone basis.
// The list of tests will be run on all zones. Different sets of tests
// requires separate requests.
type TestRequest struct {
	Zones []ZoneInfo `json:"Zones"` // List of zones to run the tests on
	Tests []string   `json:"Tests"` // List of tests to run on the zones
}

// String
//
// NOTE: Debug function.
func (t TestRequest) String() string {
	var out string
	jtr, err := json.Marshal(t)
	if err != nil {
		out = err.Error()
	} else {
		out = string(jtr)
	}

	return out
}

// ZoneTest
//
// Contains all data needed to run test as well as resulting messages
// Tries keeps track of the number of the times the test has been tried.
type ZoneTest struct {
	Status   int8 // See TestStatusMap for ref.
	Tries    int8
	Weight   int8      // Helps determine the order/grouping of tests.
	DelegNS  []DelegNS // Only used with undelegated tests
	Messages []string  // Key for message table
}

type ZoneInfo struct {
	Name    string  `json:"Name"`    // Name of the zone
	DelegNS DelegNS `json:"DelegNS"` // Delegation info (optional)
}

// DelegNS
//
// (optional) Manually supplied information about delegating name servers.
type DelegNS struct {
	Name string `json:"Name"`
	IP   string `json:"IP"`
	DS   string `json:"DS"`
}

// TestMessage
//
// Provides both a summary (Short) and a more thorough message.
type TestMessage struct {
	Short string `json:"Short"`
	Long  string `json:"Long"`
}

// TestRunner
//
// Contains test set for a zone. Tests that have passed or failed are moved to the Done queue
type TestRunner struct {
	Cfg   *cache.Config
	Zone  string
	Queue []ZoneTest
	Done  []ZoneTest
}

func NewRunner(cfg *cache.Config, zone string) TestRunner {
	var tr TestRunner
	return tr

}

type TestBatch []TestRunner

/*
func (tr *TestRunner) AddTest(test string) {

	switch test {
	case "Address01":
		tr.Queue = append(tr.Queue, Address01())
	default:
	}

}

*/

/*

Inputs:

## General
- Zone to be tested

basic01
+ test type (undelegated)
+ root name servers


address01
	list of IANA special addresses

connectivity03 + 04
    "Child Zone" - The domain name to be tested.
    "ASN Database"
    "Cymru Base Name"
    "Ris Whois Server" ASN Database is "RIPE", the default value is "riswhois.ripe.net".

consistency01
    "Child Zone" - The domain name to be tested
    "Accepted Serial Difference" - Accepted difference between SOA serial values from SOA records of different name servers for Child Zone. Default value is 0, i.e. no difference.


syntax08
    The hostnames to be tested. The hostnames comes from looking up the MX record for the domain being tested.



*/
