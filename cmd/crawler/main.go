package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/frzifus/vlookup/pkg/macpack"
	"github.com/frzifus/vlookup/pkg/version"
)

func main() {
	var (
		srcFetchMacLarge  = flag.Bool("src.fetch-l", false, "get large from ieee.org")
		srcFetchMacMedium = flag.Bool("src.fetch-m", false, "get medium from ieee.org")
		srcFetchMacSmall  = flag.Bool("src.fetch-s", false, "get small from ieee.org")
		srcFetchAll       = flag.Bool("src.fetch-all", false, "get all from ieee.org")
		srcFetchCustom    = flag.String("src.fetch-custom", "", "get small from ieee.org")

		timeout = flag.Duration("timeout", 30*time.Second, "specified timeout")

		store = flag.String("o", "", "output file")

		printVersion = flag.Bool("version", false, "print version")
	)
	flag.Parse()
	if *printVersion {
		fmt.Println(version.Version())
		return
	}

	if *srcFetchAll {
		*srcFetchMacLarge, *srcFetchMacMedium, *srcFetchMacSmall = true, true, true
	}
	urls := crawlURLS(*srcFetchMacLarge, *srcFetchMacMedium, *srcFetchMacSmall, *srcFetchCustom)
	if len(urls) == 0 {
		flag.PrintDefaults()
		return
	}

	t := time.Now()
	c := &http.Client{Timeout: *timeout}
	for i, url := range urls {
		log.Printf("%d) get: %s", i, url)
		resp, err := c.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		name := "unknown"
		if *store != "" {
			name = *store
		}
		f, err := os.Create(fmt.Sprintf("%s_%d_%s.csv", t.Format("2006-01-02"), i, name))
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		if _, err = io.Copy(f, resp.Body); err != nil {
			log.Fatalln(err)
		}
	}
}

func crawlURLS(large, medium, small bool, custom string) []string {
	var urls []string
	if large {
		urls = append(urls, macpack.RemoteIeeeMACLarge)
	}
	if medium {
		urls = append(urls, macpack.RemoteIeeeMACMedium)
	}
	if small {
		urls = append(urls, macpack.RemoteIeeeMACSmall)
	}
	if custom != "" {
		urls = append(urls, custom)
	}
	return urls
}
