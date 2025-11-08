package zonetests

import (
	"encoding/json"
	"zonetree/cache"
)

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

// ZoneTest
//
// All tests needs to implement the interface methods.
type ZoneTest interface {
	New(cfg *cache.Config)
	Status()
	Run()
}

// TestRequest
//
// Client request to run a number of tests on a nmber of zones.
// Delegation information is optional and will only be used if.
// Undelegated flag is set to true
// The list of tests will be run on all zones. Different sets of tests
// requires separate requests.
type TestRequest struct {
	Zones       []ZoneInfo `json:"ZoneInfo"`    // List of zones to run the tests on
	Tests       []string   `json:"Tests"`       // List of tests to run on the zones
	Undelegated bool       `json:"Undelegated"` // If true, use the name server data provided
}

type ZoneInfo struct {
	Name string
	UNS  []cache.NSIP // (Undelegated) Name Servers
	DS   string
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
	//	Cfg   *cache.Config
	Queue []ZoneTest
	Done  []ZoneTest
}

func NewRunner(cfg *cache.Config, zone string) TestRunner {
	var tr TestRunner
	return tr

}

// TestBatch
//
// A set of TestRunners
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
