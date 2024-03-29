package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/frzifus/vlookup/pkg/arp"
	"github.com/frzifus/vlookup/pkg/macpack"
	"github.com/frzifus/vlookup/pkg/tables"
	"github.com/frzifus/vlookup/pkg/version"
)

const (
	format = "%-5s %-10s %-20s %-20s %-20s %-15s\n"
)

func main() {
	var (
		source = flag.String("src", "embd-l", "options: ieee-s, ieee-m, ieee-l, embd-s, embd-m, embd-l")

		srcLocalFile = flag.String("src.local-file", "", "use file input")

		trimAddress = flag.Int("trim.address", 40, "limits the length of the address field")

		arpScan    = flag.Bool("arp.scan", true, "actively searches the network for other devices, this operation requires root privileges")
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

	opts, err := srcOptions(*source, *srcLocalFile)
	if err != nil {
		log.Fatalln(err)
	}
	if len(opts) == 0 {
		flag.PrintDefaults()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), *arpTimeout)
	defer cancel()
	var scanResult []*arp.Entry
	if *arpScan {
		if scanResult, err = doScan(ctx, *iface); err != nil {
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
		i++
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

func srcOptions(source string, local string) ([]macpack.Option, error) {
	var opts []macpack.Option
	if local != "" {
		opts = append(opts, macpack.WithLocalSource(local))
		return opts, nil
	}

	fs := tables.Get()
	switch source {
	case "":
		return nil, errors.New("missing data source")
	case "ieee-l":
		opts = append(opts, macpack.WithRemoteSource(macpack.RemoteIeeeMACLarge))
	case "ieee-m":
		opts = append(opts, macpack.WithRemoteSource(macpack.RemoteIeeeMACLarge))
	case "ieee-s":
		opts = append(opts, macpack.WithRemoteSource(macpack.RemoteIeeeMACSmall))
	case "embd-l":
		f, err := fs.Open(path.Base(macpack.RemoteIeeeMACLarge))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		opts = append(opts, macpack.WithReaderSource(f))
	case "embd-m":
		f, err := fs.Open(path.Base(macpack.RemoteIeeeMACMedium))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		opts = append(opts, macpack.WithReaderSource(f))
	case "embd-s":
		f, err := fs.Open(path.Base(macpack.RemoteIeeeMACSmall))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		opts = append(opts, macpack.WithReaderSource(f))
	default:
		opts = append(opts, macpack.WithLocalSource(local))
	}
	return opts, nil
}

func doScan(ctx context.Context, use string) ([]*arp.Entry, error) {
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
		if use != "" && use != iface.Name {
			continue
		}

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
