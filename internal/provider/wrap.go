package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

func WrapRunMethod(ctx context.Context, kctx *kong.Context) error {
	for _, node := range kctx.Model.Children {
		if node.Type == kong.CommandNode {
			inter := node.Target.Interface()
			intertype := reflect.TypeOf(inter)
			zerolog.Ctx(ctx).Debug().Str("type", intertype.String()).Msg("Wrapping Run method")
		}
	}
	return nil
}

// returns the command[I] that kong has selected
func GetSelectedCommand(_ context.Context, kctx *kong.Context) (CommandRunner, error) {

	curr := kctx.Selected()

	if curr == nil {
		return nil, fmt.Errorf("no command selected")
	}

	if curr.Target.Kind() != reflect.Struct {
		return nil, fmt.Errorf("command target is not a struct")
	}

	if cr, ok := curr.Target.Addr().Interface().(CommandRunner); ok {
		return cr, nil
	}

	return nil, fmt.Errorf("command target is not a CommandRunner")
}

func RunSelectedCommand(ctx context.Context, kctx *kong.Context, cp ContentProvider) error {

	ctx = BindToKongContext(ctx, kctx)

	cmd, err := GetSelectedCommand(ctx, kctx)
	if err != nil {
		return err
	}

	hlpr := cmd.AnyHelper()

	_, err = hlpr.Start(ctx, cp)
	return err
}
