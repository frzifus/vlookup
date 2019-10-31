package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/frzifus/vlookup/pkg/arp"
	"github.com/frzifus/vlookup/pkg/macpack"
	"github.com/frzifus/vlookup/pkg/version"
)

const (
	format = "%-5s %-10s %-20s %-20s %-20s %-15s\n"
)

func main() {
	var (
		srcFetchMacLarge  = flag.Bool("src.fetch-l", false, "get large from ieee.org")
		srcFetchMacMedium = flag.Bool("src.fetch-m", false, "get medium from ieee.org")
		srcFetchMacSmall  = flag.Bool("src.fetch-s", false, "get small from ieee.org")
		srcLocalFile      = flag.String("src.local-file", "", "use file input")

		trimAddress = flag.Int("trim.address", 40, "limits the length of the address field")

		iface = flag.String("i", "", "filter interface")
		store = flag.String("o", "", "output file")

		printVersion = flag.Bool("version", false, "print version")
	)
	flag.Parse()
	if *printVersion {
		fmt.Println(version.Version())
		return
	}

	opts := srcOptions(*srcFetchMacLarge, *srcFetchMacMedium, *srcFetchMacSmall, *srcLocalFile)
	if len(opts) == 0 {
		flag.PrintDefaults()
		return
	}
	mp, err := macpack.New(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("check %d entries\n", len(mp))

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, "idx", "interface", "IP", "MAC", "Name", "Address")
	fmt.Fprintf(&buf, format, "---", "---------", "--", "---", "----", "-------")
	for i, e := range arp.ParseEntries(arp.FromCache()) {
		if *iface != "" && e.Device != nil && e.Device.Name != *iface {
			continue
		}
		devIface := "unknown"
		if e.Device != nil {
			devIface = e.Device.Name
		}
		idx, mac := strconv.Itoa(i), e.Mac.String()
		name, addr := "not found", ""
		if o := mp.Get(mac); o != nil {
			name, addr = o.Name, o.Address
			if len(addr) > *trimAddress {
				addr = addr[0:*trimAddress]
			}
		}
		ip := e.Address.String()
		fmt.Fprintf(&buf, format, idx, devIface, ip, mac, name, addr)
	}
	var b io.Reader = &buf
	if *store != "" {
		f, err := os.Create(*store)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		b = io.TeeReader(b, f)
	}
	if _, err := io.Copy(os.Stdout, b); err != nil {
		log.Fatalln(err)
	}
}

func srcOptions(large, medium, small bool, local string) []macpack.Option {
	var opts []macpack.Option
	if large {
		opts = append(opts, macpack.WithRemoteSource(macpack.RemoteIeeeMACLarge))
	}
	if medium {
		opts = append(opts, macpack.WithRemoteSource(macpack.RemoteIeeeMACLarge))
	}
	if small {
		opts = append(opts, macpack.WithRemoteSource(macpack.RemoteIeeeMACSmall))
	}
	if local != "" {
		opts = append(opts, macpack.WithLocalSource(local))
	}
	return opts
}
