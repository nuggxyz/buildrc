package full

import (
	"context"
	"encoding/json"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/buildrc/pkg/buildrc"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"
	"github.com/walteh/terrors"
)

var _ snake.Flagged = (*Handler)(nil)

type Handler struct {
	FilesDir string `json:"files-dir"`
}

func (me *Handler) Flags(flgs *pflag.FlagSet) {
	flgs.StringVarP(&me.FilesDir, "files-dir", "f", "", "Write all files and buildrc.json to this directory")
}

func (me *Handler) Cobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "full",
		Short: "builds buildrc metadata files",
	}

	return cmd
}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, fls afero.Fs) error {

	gitp, err := git.NewGitGoGitProvider(fls, ".")
	if err != nil {
		return err
	}

	revision, err := buildrc.GetBuildrcJSON(ctx, gitp)
	if err != nil {
		return err
	}

	byt, err := json.Marshal(revision)
	if err != nil {
		return err
	}

	if me.FilesDir != "" {
		mapper, err := revision.Files()
		if err != nil {
			return err
		}

		raw := afero.NewOsFs()

		fs := afero.NewBasePathFs(raw, me.FilesDir)

		err = fs.MkdirAll(me.FilesDir, 0755)
		if err != nil {
			return err
		}

		for k, v := range mapper {
			err = afero.WriteFile(fs, k, []byte(v), 0644)
			if err != nil {
				return terrors.Wrap(err, "unable to write file")
			}
		}

		err = afero.WriteFile(fs, "buildrc.json", byt, 0644)
		if err != nil {
			return err
		}
	}

	cmd.Printf("%s\n", string(byt))

	return nil
}
