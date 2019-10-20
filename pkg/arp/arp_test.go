package arp

import (
	"bytes"
	"io"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseEntries(t *testing.T) {
	mac1, err := net.ParseMAC("aa:bb:cc:dd:ee:ff")
	if err != nil {
		t.Fatal(err)
	}
	mac2, err := net.ParseMAC("ff:ee:dd:cc:bb:aa")
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		name string
		r    io.Reader
		want []*Entry
	}{
		{
			name: "expected",
			r: bytes.NewBuffer([]byte(
				`IP address       HW type     Flags       HW address            Mask     Device
				 192.168.1.1      0x1         0x2         aa:bb:cc:dd:ee:ff     *        unknown
				 192.168.1.2      0x1         0x2         ff:ee:dd:cc:bb:aa     *        unknown
				`)),
			want: []*Entry{
				{
					Address: net.ParseIP("192.168.1.1"),
					Mac:     mac1,
					Mask:    "*",
				},
				{
					Address: net.ParseIP("192.168.1.2"),
					Mac:     mac2,
					Mask:    "*",
				},
			},
		},
		{
			name: "empty cache",
			r: bytes.NewBuffer([]byte(
				`IP address       HW type     Flags       HW address            Mask     Device
				`)),
			want: []*Entry{},
		},
		{
			name: "no table",
			r:    bytes.NewBuffer([]byte{}),
			want: []*Entry{},
		},
		{
			name: "invalid ip or mac",
			r: bytes.NewBuffer([]byte(
				`IP address       HW type     Flags       HW address            Mask     Device
				 xxxxxxxxxxx      0x1         0x2         aa:bb:cc:dd:ee:ff     *        unknown
				 192.168.1.2      0x1         0x2         xxxxxxxxxxxxxxxxx     *        unknown
				`)),
			want: []*Entry{
				{
					Mac:  mac1,
					Mask: "*",
				},
				{
					Address: net.ParseIP("192.168.1.2"),
					Mask:    "*",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := ParseEntries(tc.r); !cmp.Equal(got, tc.want) {
				t.Error(cmp.Diff(got, tc.want))
			}
		})
	}
}
