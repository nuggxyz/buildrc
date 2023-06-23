package docker

import (
	"context"
)

func Run(ctx context.Context) error {
	return nil

	// st := llb.Image("docker.io/library/alpine:latest").
	// 	Run(llb.Shlex("echo 'Hello, world!' > /hello")).Root()

	// def, err := st.Marshal(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	// // Define the options for the build
	// opt := client.SolveOpt{
	// 	Exports: []client.ExportEntry{
	// 		{
	// 			Type: "image",
	// 			Attrs: map[string]string{
	// 				"name": "docker.io/username/myimage:latest",
	// 				"push": "true",
	// 			},
	// 		},
	// 	},
	// }

	// // r, err := os.UserHomeDir()
	// // if err != nil {
	// // 	return err
	// // }

	// cli, err := client.New(ctx, "tcp://0.0.0.0:1234", client.WithFailFast())
	// if err != nil {
	// 	return fmt.Errorf("failed to create client: %s", err)
	// }
	// defer cli.Close()

	// ch := make(chan *client.SolveStatus)
	// eg, ctx := errgroup.WithContext(ctx)
	// eg.Go(func() error {
	// 	res, err := cli.Solve(ctx, def, opt, ch)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	zerolog.Ctx(ctx).Debug().Any("result", res).Msg("Solve result")

	// 	return err
	// })
	// eg.Go(func() error {
	// 	var c console.Console
	// 	// Avoid getting the console of a container, but do something useful with
	// 	// the solve status.
	// 	_, err := progressui.DisplaySolveStatus(ctx, "", c, os.Stdout, ch)
	// 	return err
	// })
	// if err := eg.Wait(); err != nil {
	// 	return fmt.Errorf("failed to solve: %s", err)
	// }

	// return nil

}
