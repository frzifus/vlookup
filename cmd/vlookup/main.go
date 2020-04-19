package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

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

		arpScan    = flag.Bool("arp.scan", false, "actively searches the network for other devices, this operation requires root privileges")
		arpTimeout = flag.Duration("arp.timeout", 10*time.Second, "time to wait for responses")

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

	ctx, cancel := context.WithTimeout(context.Background(), *arpTimeout)
	defer cancel()
	var scanResult []*arp.Entry
	var err error
	if *arpScan {
		if scanResult, err = doScan(ctx); err != nil {
			log.Fatalln(err)
		}
		log.Println("finished scan")
	}

	// NOTE: to improve performance, the comparison list should be updated in
	// parallel with the network scan. Also the same context can be used for this.
	mp, err := macpack.New(opts...)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("check %d vendor entries\n", len(mp))

	// NOTE: the cache list and the scan result are merged here. Duplicates
	// are removed. In principle, this should be performed by the arp discovery
	// service.
	// TODO: move and hide in arp package
	entries := make(map[string]*arp.Entry)
	for _, e := range arp.ParseEntries(arp.FromCache()) {
		entries[e.Address.String()] = e
	}
	for _, e := range scanResult {
		entries[e.Address.String()] = e
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, "idx", "interface", "IP", "MAC", "Name", "Address")
	fmt.Fprintf(&buf, format, "---", "---------", "--", "---", "----", "-------")
	var i int
	for _, e := range entries {
		i++
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

func doScan(ctx context.Context) ([]*arp.Entry, error) {
	if os.Geteuid() > 0 {
		log.Fatalln("user has insufficient permissions")
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	hosts := make(chan arp.Entry)
	errc := make(chan error)
	var entries []*arp.Entry
	for _, iface := range ifaces {
		go func(ctx context.Context, iface net.Interface) {
			if iface.Flags&(net.FlagLoopback|net.FlagPointToPoint) != 0 ||
				iface.Flags&net.FlagUp == 0 {
				log.Println("skip interface: ", iface.Name)
				return
			}

			log.Println("start scan on interface", iface.Name)
			d, err := arp.NewDiscovery(&iface)
			if err != nil {
				errc <- err
				return
			}
			defer d.Close()
			if err := d.Find(ctx, hosts); err != nil {
				errc <- err
			}
		}(ctx, iface)
	}
	for {
		select {
		case h := <-hosts:
			entries = append(entries, &h)
		case <-ctx.Done():
			return entries, nil
		case err := <-errc:
			return nil, err
		}
	}
}
