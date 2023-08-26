package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/walteh/buildrc/cmd/root"
	"github.com/walteh/snake"
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

func main() {

	ctx := context.Background()

	rootCmd := snake.NewRootCommand(ctx, &root.Root{})

	if err := snake.DecorateRootCommand(ctx, rootCmd, &snake.DecorateOptions{
		Headings: color.New(color.FgCyan, color.Bold),
		ExecName: color.New(color.FgHiGreen, color.Bold),
		Commands: color.New(color.FgHiRed, color.Faint),
	}); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "[%s] (error) %+v\n", rootCmd.Name(), err)
		if err != nil {
			panic(err)
		}
	}

	rootCmd.SilenceErrors = true

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		_, err = fmt.Fprintf(os.Stderr, "[%s] (error) %+v\n", rootCmd.Name(), err)
		if err != nil {
			panic(err)
		}
	}

}
