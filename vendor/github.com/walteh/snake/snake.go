package snake

import (
	"context"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Snake struct {
	bindings  map[string]*reflect.Value
	resolvers map[string]Method
	root      *cobra.Command
	runlock   sync.Mutex
}

type Flagged interface {
	Flags(*pflag.FlagSet)
}

type Cobrad interface {
	Cobra() *cobra.Command
}

type NewSnakeOpts struct {
	Root                       *cobra.Command
	Commands                   []Method
	Resolvers                  []Method
	GlobalContextResolverFlags bool
}

func attachMethod(me *Snake, cmd *cobra.Command, name string, globalFlags *pflag.FlagSet) (*cobra.Command, error) {

	if cmd == nil {
		return nil, nil
	}

	if flgs, err := FlagsFor(name, func(s string) Method {
		return me.resolvers[s]
	}); err != nil {
		return nil, err
	} else {
		flgs.VisitAll(func(f *pflag.Flag) {
			if globalFlags != nil && globalFlags.Lookup(f.Name) != nil {
				return
			}
			cmd.Flags().AddFlag(f)
		})
	}

	oldRunE := cmd.RunE

	// if a flag is not set, we check the environment for "cmd_name_arg_name"

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				return
			}
			val := strings.ToUpper(me.root.Name() + "_" + strings.ReplaceAll(f.Name, "-", "_"))
			envvar := os.Getenv(val)
			if envvar == "" {
				return
			}
			err := f.Value.Set(envvar)
			if err != nil {
				return
			}
		})
		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		defer setBindingWithLock(me, cmd)()
		defer setBindingWithLock(me, args)()

		err := runResolvingArguments(name, func(s string) IsRunnable {
			return me.resolvers[s]
		}, me.bindings)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		if oldRunE != nil {
			err := oldRunE(cmd, args)
			if err != nil {
				return HandleErrorByPrintingToConsole(cmd, err)
			}
		}
		return nil
	}

	return cmd, nil

}

func NewSnake(opts *NewSnakeOpts) (*cobra.Command, error) {

	root := opts.Root

	if root == nil {
		root = &cobra.Command{}
	}

	snk := &Snake{
		bindings:  make(map[string]*reflect.Value),
		resolvers: make(map[string]Method),
		root:      root,
	}

	for _, v := range opts.Resolvers {
		for _, n := range v.Names() {
			snk.resolvers[n] = v
		}

		if opts.GlobalContextResolverFlags && v.IsContextResolver() {
			v.Flags(root.PersistentFlags())
		}
	}

	for _, v := range opts.Commands {
		for _, n := range v.Names() {
			snk.resolvers[n] = v
		}
	}

	// these will always be overwritten in the RunE function
	snk.resolvers["*cobra.Command"] = NewArgumentMethod[*cobra.Command](&inlineResolver[*cobra.Command]{
		flagFunc: func(*pflag.FlagSet) {},
		runFunc: func() (*cobra.Command, error) {
			return &cobra.Command{}, nil
		},
	})

	snk.resolvers["[]string"] = NewArgumentMethod[[]string](&inlineResolver[[]string]{
		flagFunc: func(*pflag.FlagSet) {},
		runFunc: func() ([]string, error) {
			return []string{}, nil
		},
	})

	for _, exer := range snk.resolvers {
		if exer.Command() == nil {
			continue
		}
		name := exer.Names()[0]
		if cmd, err := attachMethod(snk, exer.Command().Cobra(), name, root.PersistentFlags()); err != nil {
			return nil, err
		} else if cmd != nil {
			err := exer.ValidateResponse()
			if err != nil {
				return nil, err
			}
			root.AddCommand(cmd)
		}
	}

	if opts.GlobalContextResolverFlags {
		// this will force the context to be resolved before any command is run
		snk.resolvers["root"] = NewCommandMethod(&fakeCobraWithContext{})
	} else {
		snk.resolvers["root"] = NewCommandMethod(&fakeCobra{})
	}

	root.RunE = func(cmd *cobra.Command, args []string) error {
		err := runResolvingArguments("root", func(s string) IsRunnable {
			return snk.resolvers[s]
		}, snk.bindings)
		if err != nil {
			return HandleErrorByPrintingToConsole(cmd, err)
		}
		return nil
	}

	root.SilenceUsage = true

	return root, nil
}

func NewCommandMethod[I Cobrad](cbra I) Method {

	ec := &method{
		flags:              func(*pflag.FlagSet) {},
		validationStrategy: commandResponseValidationStrategy,
		responseStrategy:   commandResponseHandleStrategy,
		names:              []string{reflect.TypeOf((*I)(nil)).Elem().String()},
		method:             getRunMethod(cbra),
		cmd:                cbra,
	}

	if flg, ok := any(cbra).(Flagged); ok {
		ec.flags = flg.Flags
	}

	return ec
}

func NewArgumentMethod[A any](m Flagged) Method {

	ec := &method{
		flags:              m.Flags,
		validationStrategy: validate1ArgumentResponse[A],
		responseStrategy:   handle1ArgumentResponse[A],
		names:              namesBuilder((*A)(nil)),
		method:             getRunMethod(m),
	}

	return ec
}

func New2ArgumentMethod[A any, B any](m Flagged) Method {

	ec := &method{
		flags:              m.Flags,
		validationStrategy: validate2ArgumentResponse[A, B],
		responseStrategy:   handle2ArgumentResponse[A, B],
		names:              namesBuilder((*A)(nil), (*B)(nil)),
		method:             getRunMethod(m),
	}

	return ec
}

func New3ArgumentMethod[A any, B any, C any](m Flagged) Method {

	ec := &method{
		flags:              m.Flags,
		validationStrategy: validate3ArgumentResponse[A, B, C],
		responseStrategy:   handle3ArgumentResponse[A, B, C],
		names:              namesBuilder((*A)(nil), (*B)(nil), (*C)(nil)),
		method:             getRunMethod(m),
	}

	return ec
}

func namesBuilder(inter ...any) []string {

	var names []string
	for _, v := range inter {
		names = append(names, reflect.TypeOf(v).Elem().String())
	}
	return names
}

type fakeCobra struct {
}

func (me *fakeCobra) Cobra() *cobra.Command {
	return &cobra.Command{}
}

func (me *fakeCobra) Run(cmd *cobra.Command) error {
	return errors.Errorf("no method found for %q", cmd.Name())
}

type fakeCobraWithContext struct {
	internal fakeCobra
}

func (me *fakeCobraWithContext) Cobra() *cobra.Command {
	return me.internal.Cobra()
}

func (me *fakeCobraWithContext) Run(_ context.Context, cmd *cobra.Command) error {
	return me.internal.Run(cmd)
}
