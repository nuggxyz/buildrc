package full

import (
	"context"
	"encoding/json"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/buildrc"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
	FilesDir string `json:"files-dir"`
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "get current revision",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.Flags().StringVarP(&me.FilesDir, "files-dir", "", "", "The directory to write the files to")

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider, fls afero.Fs) error {

	revision, err := buildrc.GetBuildrcJSON(ctx, gitp, nil)
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

		fs := afero.NewBasePathFs(fls, me.FilesDir)

		err = fs.MkdirAll(me.FilesDir, 0755)
		if err != nil {

			return err
		}

		for k, v := range mapper {
			err = afero.WriteFile(fs, k, []byte(v), 0644)
			if err != nil {
				return err
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
