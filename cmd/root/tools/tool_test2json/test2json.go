package tool_test2json

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/test2json"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
	flagP string
	flagT bool
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "convert test output to json",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.Flags().StringVarP(&me.flagP, "pkg", "p", "", "report `pkg` as the package being tested in each event")
	cmd.Flags().BoolVarP(&me.flagT, "timestamp", "t", false, "include timestamps in events")

	return cmd
}

func (me *Handler) ParseArguments(_ context.Context, _ *cobra.Command, _ []string) error {
	return nil
}

func (me *Handler) Run(_ context.Context) error {

	var mode test2json.Mode
	if me.flagT {
		mode |= test2json.Timestamp
	}
	c := test2json.NewConverter(os.Stdout, me.flagP, mode)
	defer c.Close()

	if flag.NArg() == 0 {
		_, err := io.Copy(c, os.Stdin)
		if err != nil {
			fmt.Fprintf(c, "test2json: %v\n", err)
		}
	} else {
		args := flag.Args()
		cmd := exec.Command(args[0], args[1:]...)
		w := &countWriter{0, c}
		cmd.Stdout = w
		cmd.Stderr = w
		ignoreSignals()
		err := cmd.Run()
		if err != nil {
			if w.n > 0 {
				// Assume command printed why it failed.
			} else {
				fmt.Fprintf(c, "test2json: %v\n", err)
			}
		}
		c.Exited(err)
		if err != nil {
			err := c.Close()
			if err != nil {
				fmt.Fprintf(c, "test2json: %v\n", err)
			}
			os.Exit(1)
		}
	}

	return nil
}

// ignoreSignals ignore the interrupt signals.
func ignoreSignals() {
	signal.Ignore(signalsToIgnore...)
}

type countWriter struct {
	n int64
	w io.Writer
}

func (w *countWriter) Write(b []byte) (int, error) {
	w.n += int64(len(b))
	return w.w.Write(b)
}
