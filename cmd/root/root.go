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

	"github.com/walteh/snake"
)

func NewCommand() (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:   "retab",
		Short: "retab brings tabs to your code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return snake.NewSnake(&snake.NewSnakeOpts{
		Root: cmd,
		Resolvers: []snake.Method{
			snake.NewArgumentMethod[context.Context](&resolvers.ContextResolver{}),
			snake.NewArgumentMethod[afero.Fs](&resolvers.FSResolver{}),
			snake.NewArgumentMethod[afero.File](&resolvers.FileResolver{}),
		},
		Commands: []snake.Method{
			snake.NewCommandMethod(&full.Handler{}),
			snake.NewCommandMethod(&diff.Handler{}),
			snake.NewCommandMethod(&binary.Handler{}),
			snake.NewCommandMethod(&install.Handler{}),
			snake.NewCommandMethod(&next.Handler{}),
		},
		GlobalContextResolverFlags: true,
	})

}
