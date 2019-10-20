// +build linux

package arp

import (
	"bytes"
	"io"
	"os"
)

// FromCache returns "/proc/net/arp" as io.Reader
func FromCache() io.Reader {
	f, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil
	}
	defer f.Close()
	var b bytes.Buffer
	io.Copy(&b, f)
	return &b
}
