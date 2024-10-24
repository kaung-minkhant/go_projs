package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  fmt.Printf("domain, hasMX, hasSPF, SPFRecord, hasDMARC, DMARCRecord\n")
  for scanner.Scan() {
    checkDomain(scanner.Text())
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
}

func checkDomain(domain string) {
  var hasMX, hasSPF, hasDMARC bool = false, false, false
  var spfRecord, dmarcRecord string = "", ""

  mxRecords, err := net.LookupMX(domain)
  if err != nil {
    log.Printf("Error looking up mxRecords: %v\n", err)
  }
  if len(mxRecords) > 0 {
    hasMX = true
  }

  txt, err := net.LookupTXT(domain)
  if err != nil {
    log.Printf("Error looking up txt records: %v\n", err)
  }
  for i, record := range txt {
    fmt.Printf("Records %v: %s\n", i, record)
  }
  for _, record := range txt {
    if strings.HasPrefix(record, "v=spf1") {
      hasSPF = true
      spfRecord = record
      break
    } 
  }

  dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
  if err != nil {
    fmt.Printf("Error looking up dmarcRecords: %v\n", err)
  }

  for i, record := range dmarcRecords {
    fmt.Printf("Dmarc Records %v: %s\n", i, record)
  }
  for _, record := range dmarcRecords {
    if strings.HasPrefix(record, "v=DMARC1") {
      hasDMARC = true
      dmarcRecord = record
      break
    }
  }

  fmt.Printf("%v, %v, %v, %v, %v, %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
