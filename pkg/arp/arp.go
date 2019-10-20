package arp

import (
	"bufio"
	"encoding/hex"
	"io"
	"net"
	"strings"
)

const (
	columnIPAddr int = iota
	columnHWType
	columnFlags
	columnHWAddr
	columnMask
	columnDevice
	columnBound
)

// Entry represents an entry in the arp cache.
// This can usually be found under linux under "/proc/net/arp".
type Entry struct {
	Address net.IP
	Type    byte
	Flags   byte
	Mac     net.HardwareAddr
	Mask    string
	Device  *net.Interface
}

// ParseEntries parses s as an arp cache entry, returning the result.
// The table should look like this:
// IP address       HW type     Flags       HW address           Mask    Device
// 192.168.1.1      0x1         0x2         aa:bb:cc:dd:ee:ff    *       e.g.1
// 192.168.1.2      0x1         0x2         ff:ee:dd:cc:bb:aa    *       e.g.2
// If entries are not valid, they are ignored. If the list is empty, an empty
// result list is returned.
func ParseEntries(r io.Reader) []*Entry {
	s := bufio.NewScanner(r)
	s.Scan() // skip header
	entries := make([]*Entry, 0)
	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) < columnBound {
			continue
		}
		e := &Entry{Address: net.ParseIP(f[columnIPAddr]), Mask: f[columnMask]}
		if t, err := hex.DecodeString(f[columnHWType]); err == nil && len(t) >= 0 {
			e.Type = t[0]
		}
		if fl, err := hex.DecodeString(f[columnFlags]); err == nil && len(fl) >= 0 {
			e.Flags = fl[0]
		}
		if mac, err := net.ParseMAC(f[columnHWAddr]); err == nil {
			e.Mac = mac
		}
		if iface, err := net.InterfaceByName(f[columnDevice]); err == nil {
			e.Device = iface
		}
		entries = append(entries, e)
	}

	return entries
}
