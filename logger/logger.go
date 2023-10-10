package logger

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"runtime"
	"time"
)

var Global zerolog.Logger

var Blizzard zerolog.Logger

var Orca zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	Global = createLogger("arctic")
	Blizzard = createLogger("blizzard")
	Orca = createLogger("orca")
}

func createLogger(scope string) zerolog.Logger {
	w := zerolog.ConsoleWriter{
		Out: os.Stdout,
		FormatTimestamp: func(i interface{}) string {
			t := i.(json.Number)
			if _t, e := t.Int64(); e != nil {
				return ""
			} else {
				return "\033[0;36m " + time.UnixMilli(_t).Format("02/01/2006 15:04:05.000") + "\033[0m"
			}
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("\033[0;36m[%s] \033[0;0m%s", scope, i)
		},
	}
	return zerolog.New(w).With().Timestamp().Logger()
}

func Panic(e error, msg string, args ...interface{}) {
	if e == nil {
		return
	}
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	l := Global.Panic().Stack().Err(e)
	if ok && details != nil {
		file, line := details.FileLine(pc)
		l = l.Str("from", details.Name()).Str("line", fmt.Sprintf("%s:%d", file, line))
	}
	l.Msgf(msg, args...)
}
