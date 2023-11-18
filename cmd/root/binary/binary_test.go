package binary

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/rs/zerolog"
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
				Platform:     runtime.GOOS + "/" + runtime.GOARCH,
			}

			ctx := context.Background()

			ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().Level(zerolog.DebugLevel).WithContext(ctx)

			err := me.Run(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Binary.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
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
