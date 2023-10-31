package build

import "strconv"

var (
	Version = "dev"
	Hash    = "unknown"
	_date   = "0"
	Date    = 0
)

func init() {
	d, e := strconv.Atoi(_date)
	if e != nil {
		return
	}
	Date = d
}
