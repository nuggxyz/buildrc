package github

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/nuggxyz/buildrc/internal/action"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

const (
	CommandID = "gen_github"
)

type Handler struct {
	OutDir string `arg:"" help:"The directory to write the action.yml files to."`

	kctx *kong.Context
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider, kctx kong.Context) (err error) {
	me.kctx = &kctx
	_, err = me.Invoke(ctx, cp)
	return err
}

func (me *Handler) Invoke(ctx context.Context, cp provider.ContentProvider) (out *any, err error) {
	return provider.Wrap(CommandID, me.invoke)(ctx, cp)
}

func (me *Handler) invoke(ctx context.Context, cp provider.ContentProvider) (*any, error) {

	if me.kctx == nil {
		return nil, fmt.Errorf("kong context is nil")
	}

	for _, cmd := range me.kctx.Model.Children {
		name := cmd.Name
		help := cmd.Help

		pos := make([]action.Input, 0)
		flags := make([]action.Input, 0)

		for _, arg := range cmd.Positional {
			pos = append(pos, action.Input{
				Name:        arg.Name,
				Description: arg.Help,
				Required:    arg.Required,
				Default:     arg.Default,
			})
		}

		for _, arg := range cmd.Flags {
			flags = append(flags, action.Input{
				Name:        arg.Name,
				Description: arg.Help,
				Required:    arg.Required,
				Default:     arg.Default,
			})
		}

		cmm := action.Command{
			Name:        name,
			Description: help,
			Positional:  pos,
			Flags:       flags,
		}

		err := cmm.ToGithubAction().WriteYamlToDir(me.OutDir)
		if err != nil {
			return nil, err
		}
	}

	zerolog.Ctx(ctx).Info().Msgf("Wrote github actions to %s", me.OutDir)

	return new(interface{}), nil
}
