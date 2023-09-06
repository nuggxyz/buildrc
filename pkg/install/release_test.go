package install

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func TestDownloadGithubHttpIntegrationRelease(t *testing.T) {
	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().Level(zerolog.DebugLevel).WithContext(ctx)

	type args struct {
		fls     afero.Fs
		org     string
		name    string
		version string
		token   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"buildrc latest", args{afero.NewMemMapFs(), "walteh", "buildrc", "latest", ""}, false},
		{"gotestsum latest", args{afero.NewMemMapFs(), "gotestyourself", "gotestsum", "latest", ""}, false},
		{"buildrc v0.13.0", args{afero.NewMemMapFs(), "walteh", "buildrc", "v0.13.0", ""}, false},
		{"gotestsum v1.10.1", args{afero.NewMemMapFs(), "gotestyourself", "gotestsum", "v1.10.1", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := DownloadGithubRelease(ctx, tt.args.fls, tt.args.org, tt.args.name, tt.args.version, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("DownloadGithubRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
