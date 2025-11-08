package cache

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ipdoc
//
// Local copies of IANAs list of Special Use IP addresses as provided on their
// website.
// https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry-1.csv
// https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry-1.csv
var ipdoc = []string{
	"assets/iana-ipv6-special-registry-1.csv",
	"assets/iana-ipv4-special-registry-1.csv",
}

// SuIP
//
// Special Use IP. We only care about a few parameters here:
// The type of net (category) and whether iot is globally reachable.
type SuIP struct {
	Block  string
	Name   string
	Global bool
	Class  string
}

func DefaultSuIP() []SuIP {
	return LoadLists(ipdoc)
}

func LoadLists(files []string) []SuIP {

	var IPtable []SuIP

	for _, file := range files {
		IPtable = append(IPtable, ReadCSV(file)...)
	}

	return IPtable
}

// ReadCSV
//
// Parses CSV files and makes a classification list based on keywords.
func ReadCSV(file string) []SuIP {
	b, err := os.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	in := string(b)
	r := csv.NewReader(strings.NewReader(in))

	var IPtable []SuIP

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Print(err)
		}

		// Don't add the label row
		if record[0] != "Address Block" {
			// Turn textual true/false of field "Globally reachable" inte bool
			// Be sure to remove any comments by (split and use first value)
			boolValue, _ := strconv.ParseBool(strings.Split(record[8], "[")[0])
			// Check for multiple addresses in Address block field,
			// because IANA can't consistency...
			for bl := range strings.SplitSeq(record[0], " ") {
				if strings.Contains(bl, "/") {
					class := "Other"
					rg, _ := regexp.Compile("(Loopback|Local|This|Private)")
					if rg.MatchString(record[1]) {
						class = "Local"
					}
					rg, _ = regexp.Compile("(Documentation)")
					if rg.MatchString(record[1]) {
						class = "Documentation"
					}
					iptr := SuIP{
						Block:  strings.Trim(bl, ","),
						Name:   record[1],
						Global: boolValue,
						Class:  class,
					}
					IPtable = append(IPtable, iptr)
				}
			}
		}
	}
	return IPtable
}
