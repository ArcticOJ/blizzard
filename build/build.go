package build

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

var (
	Tag           = "dev"
	Hash          = "n/a"
	_date         = "0"
	Date    int64 = 0
	Version       = "n/a"
)

func init() {
	d, e := strconv.Atoi(_date)
	if e != nil {
		return
	}
	Date = int64(d)
	Version = Tag
	// Ignore this warning as `Tag` is set only on build time, which confuses the compiler, leading to false report.
	if Tag == "dev" {
		Version = fmt.Sprintf("%s#%s@%s", Tag, Hash, time.Unix(Date, 0).Format(time.RFC3339))
	}
	Version = fmt.Sprintf("%s with %s", Version, runtime.Version())
}
