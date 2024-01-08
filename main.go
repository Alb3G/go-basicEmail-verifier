package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var pf = fmt.Printf

func main() {
	sc := bufio.NewScanner(os.Stdin)
	pf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	for sc.Scan() {
		checkDomain(sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	mxRc, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up the MX records: %v\n", err)
	}
	if len(mxRc) > 0 {
		hasMX = true
	}

	txtRcs, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up the SPF records: %v\n", err)
	}
	for _, record := range txtRcs {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRcs, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	for _, record := range dmarcRcs {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	pf("%v, %v, %v, %v, %v, %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
