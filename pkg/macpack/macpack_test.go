package macpack

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var errFailingSource = errors.New("failed")

func withFailingSource() Option {
	return func(MacPack) error {
		return errFailingSource
	}
}

func withTestEntries(entries map[string]Organization) Option {
	return func(m MacPack) error {
		for k, v := range entries {
			m[k] = v
		}
		return nil
	}
}

func TestNew(t *testing.T) {
	tt := []struct {
		name    string
		opts    []Option
		want    MacPack
		wantErr bool
	}{
		{
			name: "expected",
			opts: []Option{
				withTestEntries(
					map[string]Organization{
						"test": Organization{Name: "test-name", Address: "East"},
					},
				)},
			want: MacPack{
				"test": Organization{Name: "test-name", Address: "East"},
			},
		},
		{
			name:    "option failed",
			opts:    []Option{withFailingSource()},
			wantErr: true,
		},
		{
			name: "empty macpack",
			opts: []Option{},
			want: MacPack{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(tc.opts...)
			if (err != nil) != tc.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !cmp.Equal(got, tc.want) {
				t.Error(cmp.Diff(got, tc.want))
			}
		})
	}
}

func TestMacPack_Get(t *testing.T) {
	tt := []struct {
		name string
		m    MacPack
		addr string
		want *Organization
	}{
		{
			name: "expected",
			m: MacPack{
				"ffffff": Organization{Name: "test-name1", Address: "East"},
				"aabbcc": Organization{Name: "test-name2", Address: "nord"},
			},
			addr: "ffffff",
			want: &Organization{Name: "test-name1", Address: "East"},
		},
		{
			name: "cutted",
			m: MacPack{
				"ffffff": Organization{Name: "test-name1", Address: "East"},
				"aabbcc": Organization{Name: "test-name2", Address: "nord"},
			},
			addr: "aabbccdd",
			want: &Organization{Name: "test-name2", Address: "nord"},
		},
		{
			name: "upper case",
			m: MacPack{
				"ffffff": Organization{Name: "test-name1", Address: "East"},
				"aabbcc": Organization{Name: "test-name2", Address: "nord"},
			},
			addr: "AABBCC",
			want: &Organization{Name: "test-name2", Address: "nord"},
		},
		{
			name: "upper case and cutted",
			m: MacPack{
				"ffffff": Organization{Name: "test-name1", Address: "East"},
				"aabbcc": Organization{Name: "test-name2", Address: "nord"},
			},
			addr: "AABBCCDDEEFF",
			want: &Organization{Name: "test-name2", Address: "nord"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.m.Get(tc.addr); !cmp.Equal(got, tc.want) {
				t.Error(cmp.Diff(got, tc.want))
			}
		})
	}
}

func Test_newMacPack(t *testing.T) {
	tt := []struct {
		name    string
		r       io.Reader
		want    MacPack
		wantErr bool
	}{
		{
			name: "expected",
			r: bytes.NewBuffer([]byte(
				`Registry,Assignment,Organization Name,Organization Address
				 MA-S,70B3D5F2F,TELEPLATFORMS,"Polbina st., 3/1 Moscow  RU 109388"
				 MA-S,70B3D5719,2M Technology,802 Greenview Drive  Grand Prairie TX US 75050`),
			),
			want: MacPack{
				"70b3d5719": {Name: "2M Technology", Address: "802 Greenview Drive  Grand Prairie TX US 75050"},
				"70b3d5f2f": {Name: "TELEPLATFORMS", Address: "Polbina st., 3/1 Moscow  RU 109388"},
			},
		},
		{
			name: "empty",
			r:    bytes.NewBuffer([]byte(`Registry,Assignment,Organization Name,Organization Address`)),
			want: MacPack{},
		},
		{
			name: "invalid line length",
			r: bytes.NewBuffer([]byte(
				`Registry,Assignment,Organization Name,Organization Address
				 MA-S,70B3D5F2F
				 MA-S,70B3D5719`),
			),
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := newMacPack(tc.r)
			if (err != nil) != tc.wantErr {
				t.Errorf("newMacPack() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !cmp.Equal(got, tc.want) {
				t.Error(cmp.Diff(got, tc.want))
			}
		})
	}
}
