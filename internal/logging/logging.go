package logging

import (
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

func init() {

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {

		short := file
		check := 0
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				check++
				if check == 3 {
					short = file[i+1:]
					break
				}
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
}

type Logger interface {
}
