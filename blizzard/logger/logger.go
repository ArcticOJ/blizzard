package logger

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var Logger zerolog.Logger

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	Logger = zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stdout,
		FormatTimestamp: func(i interface{}) string {
			t := i.(json.Number)
			if _t, e := t.Int64(); e != nil {
				return ""
			} else {
				return "\033[0;100m " + time.UnixMilli(_t).Format("02/01/2006 15:04:05") + " \033[0m"
			}
		},
	}).With().Timestamp().Logger()
}
