package zonetests

import (
	"zonetree/cache"
)

/*
type ZoneTests interface {
	Run(param TestParam)
	Passed()   //  boolean based on test criteria
	Messages() // returns list of messags from test results
}
*/

var TestStatusMap = map[int8]string{
	0: "Not Started",
	1: "Running",
	2: "Waiting for Retry",
	3: "Disabled",
	4: "Failed",
	5: "Passed",
}

type TestSuite struct {
	Zone       string       `json:"Zone"`       // Zone to run the test on
	Delegation []cache.NSIP `json:"Delegation"` // Only used with undelegated tests
	Queue      []string     `json:"Queue"`
	Results    []TestResult `json:"Results"`
}

type TestResult struct {
	Passed   bool     `json:"Passed"`
	Messages []string `json:"Messages"` // Key for message table
}

type Message struct {
	Short string `json:"Short"`
	Long  string `json:"Long"`
}

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
