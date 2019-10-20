package macpack

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	columnRegistry int = iota
	columnAssignment
	columnName
	columnAddress
	columnBound
)

// An Option configures a MacPack at creation time.
type Option func(mp MacPack) error

// WithRemoteSource adds entries from a remote location to the macpack register.
// e.g. path: http://example.com/list.csv
// The csv source should be formatted represent the following layout:
// - Registry,Assignment,Organization Name,Organization Address
// - MA-L, MA-S,AAAAAAAAA, orga1, A street Moscow RU 1234
func WithRemoteSource(path string) Option {
	return func(m MacPack) error {
		c := &http.Client{Timeout: 30 * time.Second}
		resp, err := c.Get(path)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var b bytes.Buffer
		if _, err = io.Copy(&b, resp.Body); err != nil {
			return err
		}
		mp, err := newMacPack(&b)
		if err != nil {
			return err
		}
		for k, v := range mp {
			m[k] = v
		}
		return nil
	}
}

// WithLocalSource adds entries from a local location to the macpack register.
// e.g. path: /opt/list.csv
// The csv source should be formatted represent the following layout:
// - Registry,Assignment,Organization Name,Organization Address
// - MA-L, MA-S,AAAAAAAAA, orga1, A street Moscow RU 1234
func WithLocalSource(path string) Option {
	return func(m MacPack) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		mp, err := newMacPack(f)
		if err != nil {
			return err
		}
		for k, v := range mp {
			m[k] = v
		}
		return nil
	}
}

// Organization contains name and address of organization
type Organization struct {
	Name    string
	Address string
}

// MacPack is used to store relations between the vendor part of a physical
// hardware address and the organization behind it.
type MacPack map[string]Organization

// New returns a new MacPack which can be modified with options.
// If an error occurs, it will be returned.
func New(opts ...Option) (MacPack, error) {
	m := make(MacPack)
	for _, o := range opts {
		if err := o(m); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// Get returns the organization that belongs to the given address.
// If there is no entry for the address nil will be returned.
// Upper and lower case are accepted as well as the following notations:
// - FF:FF:FF:FF:FF:FF or ffffffffffff
func (m MacPack) Get(addr string) *Organization {
	addr = strings.ReplaceAll(addr, ":", "")
	addr = strings.ToLower(addr)
	for ; len(addr) > 0; addr = addr[:len(addr)-1] {
		if o, ok := m[addr]; ok {
			return &o
		}
	}
	return nil
}

func newMacPack(r io.Reader) (MacPack, error) {
	reader := csv.NewReader(r)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	m := make(MacPack)
	for _, rec := range records {
		if len(rec) < columnBound {
			// TODO: logger.Warn(decoding, to_short)
			continue
		}
		m[strings.ToLower(rec[columnAssignment])] = Organization{
			Name:    rec[columnName],
			Address: rec[columnAddress],
		}
	}
	return m, nil
}
