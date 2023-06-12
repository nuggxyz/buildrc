package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
)

func NewVerboseLogger() *zerolog.Logger {

	consoleOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro, NoColor: false}

	consoleOutput.FormatMessage = func(i interface{}) string {
		if i == nil {
			return "nil"
		}

		return color.New(color.FgHiWhite).Sprintf("%s", i.(string))
	}

	pretty := pp.New()

	pretty.SetColorScheme(pp.ColorScheme{})

	prettyerr := pp.New()
	prettyerr.SetExportedOnly(false)

	consoleOutput.FormatFieldValue = func(i interface{}) string {

		switch i := i.(type) {
		case error:
			return prettyerr.Sprint(i)
		case []byte:
			var g interface{}
			err := json.Unmarshal(i, &g)
			if err != nil {
				return pretty.Sprint(string(i))
			} else {
				return pretty.Sprint(g)
			}
		}

		return pretty.Sprint(i)
	}

	consoleOutput.FormatTimestamp = func(i interface{}) string {
		return color.New(color.FgHiWhite).Sprintf("%s", time.Now().Format("[15:04:05.000000]"))
	}

	consoleOutput.FormatCaller = func(i interface{}) string {
		a := i.(string)
		tot := strings.Split(a, "/")
		if len(tot) == 3 {
			num := strings.Split(tot[2], ":")

			padding := " "

			return fmt.Sprintf("[%s:%s] %s:%s%s", color.New(color.FgHiWhite).Sprintf("%s", tot[0]), tot[1], color.New(color.FgHiWhite).Sprintf("%s", num[0]), color.New(color.FgHiWhite).Sprintf("%s", num[1]), padding)
		}

		return fmt.Sprintf("\x1b[0m\x1b[34;1m%s\x1b[0m", i)
	}

	consoleOutput.PartsOrder = []string{"level", "time", "caller", "message"}

	consoleOutput.FieldsExclude = []string{"handler", "tags"}

	l := zerolog.New(consoleOutput).With().Caller().Timestamp().Logger()

	return &l

}
