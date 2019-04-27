package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jamesog/iptoasn"
	"github.com/olekukonko/tablewriter"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [ <ip>... | <asn>... ]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// Positional parameters can be either an IP address or an AS number.
	// Evaluate all parameters and store them in slices to print separate
	// tables for IPs and ASes as they have different fields.
	var ips, asns []string
	for _, arg := range os.Args[1:] {
		addr := net.ParseIP(arg)
		switch {
		case addr != nil:
			ips = append(ips, arg)
		case strings.ToLower(arg[0:2]) == "as":
			asns = append(asns, arg)
		default:
			fmt.Fprintf(os.Stderr, "%s is neither IP address nor AS\n", arg)
			usage()
		}
	}

	if len(ips) > 0 {
		whoisIP(ips)
	}
	if len(asns) > 0 {
		if len(ips) > 0 {
			fmt.Println()
		}
		whoisAS(asns)
	}
}

func minimalTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	return table
}

func whoisIP(ips []string) {
	origins := make([]iptoasn.IP, 0, len(ips))
	for _, ip := range ips {
		ipinfo, err := iptoasn.LookupIP(ip)
		if err != nil {
			log.Fatal(err)
		}
		origins = append(origins, ipinfo)
	}
	table := minimalTable()
	table.SetHeader([]string{"AS", "IP", "BGP Prefix", "CC", "Registry", "Allocated", "AS Name"})
	for _, ip := range origins {
		as := strconv.FormatUint(uint64(ip.ASNum), 10)
		table.Append([]string{as, ip.IP, ip.BGPPrefix, ip.Country, ip.Registry, ip.Allocated, ip.ASName})
	}
	table.Render()
}

func whoisAS(asns []string) {
	origins := make([]iptoasn.ASN, 0, len(asns))
	for _, as := range asns {
		ipinfo, err := iptoasn.LookupASN(as)
		if err != nil {
			log.Fatal(err)
		}
		origins = append(origins, ipinfo)
	}
	table := minimalTable()
	table.SetHeader([]string{"AS", "CC", "Registry", "Allocated", "AS Name"})
	for _, ip := range origins {
		as := strconv.FormatUint(uint64(ip.ASNum), 10)
		table.Append([]string{as, ip.Country, ip.Registry, ip.Allocated, ip.ASName})
	}
	table.Render()
}
