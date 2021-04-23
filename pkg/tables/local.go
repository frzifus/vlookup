package	tables

import "embed"

//go:embed *.csv
var f embed.FS

// Get returns the embedded lookup tables
func Get() embed.FS {
	return f
}
