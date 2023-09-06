package binary

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// go run ./cmd binary --repository=gotestsum --organization=gotestyourself --outfile=./bin/gotestsum-binary --debug

func TestBinaryHttpIntegrationWithGithub(t *testing.T) {

	type args struct {
		Organization string
		Repository   string
		Version      string
		Token        string
		Provider     string

		versionCmd string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "gotestsum latest",
			args: args{
				Organization: "gotestyourself",
				Repository:   "gotestsum",
				Version:      "latest",
				Token:        "",
				Provider:     "github",
				versionCmd:   "--version",
			},
			wantErr: false,
		},
		{
			name: "gotestsum v1.10.1",
			args: args{
				Organization: "gotestyourself",
				Repository:   "gotestsum",
				Version:      "v1.10.1",
				Token:        "",
				Provider:     "github",
				versionCmd:   "--version",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dir := os.TempDir()

			me := Handler{
				Organization: tt.args.Organization,
				Repository:   tt.args.Repository,
				Version:      tt.args.Version,
				Token:        tt.args.Token,
				Provider:     tt.args.Provider,
				OutFile:      filepath.Join(dir, tt.args.Repository+"-binary-for-test"),
			}

			ctx := context.Background()

			ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().Level(zerolog.DebugLevel).WithContext(ctx)

			cmd := cobra.Command{}

			err := me.ParseArguments(ctx, &cmd, []string{})

			if (err != nil) != tt.wantErr {
				t.Errorf("Binary.ParseArguments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = me.Run(ctx, &cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Binary.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			defer func() {
				err := os.Remove(me.OutFile)
				if err != nil {
					t.Errorf("Error removing file: %v", err)
				}
			}()

			// try to run the file as an executable
			err = exec.CommandContext(ctx, me.OutFile, tt.args.versionCmd).Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Binary.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}

}
