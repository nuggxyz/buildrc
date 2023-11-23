package root

import (
	"context"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/cmd/root/binary"
	"github.com/walteh/buildrc/cmd/root/diff"
	"github.com/walteh/buildrc/cmd/root/full"
	"github.com/walteh/buildrc/cmd/root/install"
	"github.com/walteh/buildrc/cmd/root/next"
	"github.com/walteh/buildrc/cmd/root/resolvers"
	"github.com/walteh/buildrc/cmd/root/revision"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/simver"

	"github.com/walteh/snake"
)

func NewCommand() (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:   "buildrc",
		Short: "build time metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return snake.NewSnake(&snake.NewSnakeOpts{
		Root: cmd,
		Resolvers: []snake.Method{
			snake.NewArgumentMethod[context.Context](&resolvers.ContextResolver{}),
			snake.NewArgumentMethod[afero.Fs](&resolvers.FSResolver{}),
			snake.NewArgumentMethod[git.GitProvider](&resolvers.GitResolver{}),
			snake.New3ArgumentMethod[simver.Execution, simver.GitProvider, simver.TagReader](&resolvers.SimverResolver{}),
		},
		Commands: []snake.Method{
			snake.NewCommandMethod(&full.Handler{}),
			snake.NewCommandMethod(&diff.Handler{}),
			snake.NewCommandMethod(&binary.Handler{}),
			snake.NewCommandMethod(&install.Handler{}),
			snake.NewCommandMethod(&next.Handler{}),
			snake.NewCommandMethod(&revision.Handler{}),
		},
		GlobalContextResolverFlags: true,
	})

}
