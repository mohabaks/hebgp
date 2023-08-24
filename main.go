// simple tool to query for information about the IP, Subnet and ASN from
// https://bgp.he.net
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// BaseURL is the base URL of the BGP website
const BaseURL = "https://bgp.he.net"

// IPInfo represents information about an IP address
type IPInfo struct {
	ASN         string `json:"asn"`
	Network     string `json:"network"`
	Description string `json:"description"`
}

// NETInfo represents information about a network block
type NETInfo struct {
	ASN         string `json:"asn"`
	Network     string `json:"network"`
	Description string `json:"description"`
}

// ASNInfo represents information about an ASN number
type ASNInfo struct {
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
}

// ORGInfo represents information about an organization
type ORGInfo struct {
	Result      string `json:"result"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func main() {
	// Initialize command-line parameters
	getASN := flag.String("asn", "", "Query for ASN")
	getIP := flag.String("ip", "", "Query for IP")
	getNET := flag.String("net", "", "Query for network block")
	getORG := flag.String("org", "", "Query for organization")
	getHelp := flag.Bool("h", false, "Show help message")
	flag.Parse()

	// Show help message
	if len(os.Args[1:]) == 0 || *getHelp {
		showHelpMessage()
	}

	// Query for ASN information
	if *getASN != "" {
		url := fmt.Sprintf("%s/%s", BaseURL, *getASN)
		queryAndPrint(url, queryASN)
	}

	// Query for IP information
	if *getIP != "" {
		url := fmt.Sprintf("%s/ip/%s", BaseURL, *getIP)
		queryAndPrint(url, queryIP)
	}

	// Query for network block information
	if *getNET != "" {
		url := fmt.Sprintf("%s/net/%s", BaseURL, *getNET)
		queryAndPrint(url, queryNET)
	}

	// Query for organization information
	if *getORG != "" {
		url := fmt.Sprintf("%s/search?search[search]=%s&commit=Search",
			BaseURL, *getORG)
		queryAndPrint(url, queryORG)
	}
}

// queryAndPrint takes a url and a query function that get passed to queryParser
// for further processing.
func queryAndPrint(url string, queryFunc func(*goquery.Document)) {
	doc := queryParser(url)
	queryFunc(doc)
}

// queryParser queries a URL, parses the HTML document using goquery, and returns
// the document for further processing.
func queryParser(url string) *goquery.Document {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// check for status code error
	if res.StatusCode != 200 {
		log.Printf("status code error: %d", res.StatusCode)
	}

	// load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// queryIP query for information about the IP address and print json results
func queryIP(doc *goquery.Document) {
	var rows []IPInfo

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		asn := strings.TrimSpace(row.Find("td").Eq(0).Text())
		net := strings.TrimSpace(row.Find("td").Eq(1).Text())
		des := strings.TrimSpace(row.Find("td").Eq(2).Text())

		res := IPInfo{ASN: asn, Network: net, Description: des}
		rows = append(rows, res)
	})

	printJSON(rows)
}

// queryNET query for Network Address block and print json results
func queryNET(doc *goquery.Document) {
	var rows []NETInfo

	doc.Find("#netinfo tbody tr").Each(func(i int, row *goquery.Selection) {
		asn := strings.TrimSpace(row.Find("td").Eq(0).Text())
		net := strings.TrimSpace(row.Find("td").Eq(1).Text())
		des := strings.TrimSpace(row.Find("td").Eq(2).Text())

		res := NETInfo{ASN: asn, Network: net, Description: des}
		rows = append(rows, res)

	})

	printJSON(rows)
}

// queryORG query for network information using organization name and print
// results in json
func queryORG(doc *goquery.Document) {
	var rows []ORGInfo

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		result := strings.TrimSpace(row.Find("td").Eq(0).Text())
		kind := strings.TrimSpace(row.Find("td").Eq(1).Text())
		des := strings.TrimSpace(row.Find("td").Eq(2).Text())

		res := ORGInfo{Result: result, Type: kind, Description: des}
		rows = append(rows, res)

	})

	printJSON(rows)
}

// queryASN query for ASN number and print results in json
func queryASN(doc *goquery.Document) {
	var rows []ASNInfo

	doc.Find("#table_prefixes4 tbody tr").Each(func(i int,
		row *goquery.Selection) {
		pref := strings.TrimSpace(row.Find("td").Eq(0).Text())
		des := strings.TrimSpace(row.Find("td").Eq(1).Text())

		res := ASNInfo{Prefix: pref, Description: des}
		rows = append(rows, res)
	})

	printJSON(rows)
}

// printJSON Print the given data as JSON
func printJSON(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(jsonData))
}

// showHelpMessage print the help message
func showHelpMessage() {
	fmt.Printf("Usage: %s [OPTIONS]\n\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nExamples:")
	fmt.Printf("\n  %s -asn AS63293", os.Args[0])
	fmt.Printf("\n  %s -ip 1.1.1.1", os.Args[0])
	fmt.Printf("\n  %s -net 41.223.111.0/22", os.Args[0])
	fmt.Printf("\n  %s -org facebook\n", os.Args[0])
}
