/*
 * NETCAP - Traffic Analysis Framework
 * Copyright (c) 2017 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/dreadl0ck/netcap"
	"github.com/dreadl0ck/netcap/types"
	"github.com/mgutz/ansi"
)

func printHeader() {
	netcap.PrintLogo()
	fmt.Println()
	fmt.Println("usage examples:")
	fmt.Println("	$ net.cap -r dump.pcap")
	fmt.Println("	$ net.cap -iface eth0")
	fmt.Println("	$ net.cap -r TCP.ncap.gz")
	fmt.Println("	$ net.cap -fields -r TCP.ncap.gz")
	fmt.Println("	$ net.cap -r TCP.ncap.gz -select Timestamp,SrcPort,DstPort > tcp.csv")
	fmt.Println()
}

// usage prints the use
func printUsage() {
	printHeader()
	flag.PrintDefaults()
}

// CheckFields checks if the separator occurs inside fields of audit records
// to prevent this breaking the generated CSV file
func checkFields() {

	r, err := netcap.Open(*flagInput)
	if err != nil {
		panic(err)
	}

	var (
		h                 = r.ReadHeader()
		record            = netcap.InitRecord(h.Type)
		numExpectedFields int
	)
	if p, ok := record.(types.AuditRecord); ok {
		numExpectedFields = len(p.CSVHeader())
	} else {
		log.Fatal("netcap type does not implement the types.AuditRecord interface!")
	}

	// for {
	// 	err = r.Next(record)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		break
	// 	}

	// 	if p, ok := record.(types.AuditRecord); ok {
	// 		fields := p.CSVRecord()
	// 		// TODO refactor to use netcap lib to read file instead of calling it as command
	// 	}
	// }

	r.Close()

	// call netcap and parse output line by line
	out, err := exec.Command("netcap", "-r", *flagInput).Output()
	if err != nil {
		panic(err)
	}

	// iterate over lines
	for _, line := range strings.Split(string(out), "\n") {
		count := strings.Count(line, *flagSeparator)
		if count != numExpectedFields-1 {
			fmt.Println(strings.Replace(line, *flagSeparator, ansi.Red+*flagSeparator+ansi.Reset, -1), ansi.Red, count, ansi.Reset)
		}
	}
}
