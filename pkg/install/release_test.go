package install

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func TestInstallLatestGithubRelease(t *testing.T) {
	ctx := context.Background()

	ctx = zerolog.New(zerolog.NewConsoleWriter()).With().Caller().Logger().Level(zerolog.DebugLevel).WithContext(ctx)

	type args struct {
		fls   afero.Fs
		ofs   afero.Fs
		org   string
		name  string
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"a", args{afero.NewMemMapFs(), afero.NewMemMapFs(), "walteh", "buildrc", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallLatestGithubRelease(ctx, tt.args.ofs, tt.args.fls, tt.args.org, tt.args.name, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("InstallLatestGithubRelease() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				hd, err := os.UserHomeDir()
				if err != nil {
					t.Fatalf("Error getting home dir: %v", err)
				}
				_, err = tt.args.ofs.Open(filepath.Join(hd, "."+tt.args.name, tt.args.name))
				if err != nil {
					t.Fatalf("Error opening file: %v", err)
				}
			}

		})
	}
}
