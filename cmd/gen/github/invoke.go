package github

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/internal/action"
	"github.com/nuggxyz/buildrc/internal/provider"
)

type Handler struct {
	OutDir string `arg:"" help:"The directory to write the action.yml files to."`
}

func (me *Handler) Invoke(ctx context.Context, cp provider.ContentProvider) (*OUTPUT, error) {

	k := provider.GetKongContext(ctx)

	for _, cmd := range k.Model.Children {
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

	fmt.Println("action.yml files written to", me.OutDir)

	return new(interface{}), nil
}
