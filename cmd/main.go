package main

import (
	"context"
	"log"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/walteh/buildrc/cmd/root"
	"github.com/walteh/buildrc/internal/git"
	"github.com/walteh/buildrc/internal/github/actions"
	"github.com/walteh/buildrc/internal/github/restapi"
	"github.com/walteh/buildrc/internal/pipeline"
	"github.com/walteh/snake"
)

func main() {

	ctx := context.Background()

	execgit := git.NewGitGoGitProvider()

	var pr git.PullRequestProvider
	var release git.ReleaseProvider
	var repometa git.RemoteRepositoryMetadataProvider
	var pipe pipeline.Pipeline
	var fs afero.Fs

	if actions.IAmInAGithubAction(ctx) {
		actionpipe, err := actions.NewGithubActionPipeline(ctx)
		if err != nil {

			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create content provider")
			return
		}

		pipe = actionpipe

		ghrestapi, err := restapi.NewGithubClient(ctx, execgit, actionpipe)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create github client")
			return
		}

		pr = ghrestapi
		release = ghrestapi
		repometa = ghrestapi

		fs = afero.NewOsFs()
	} else {
		zerolog.Ctx(ctx).Warn().Msg("not running in github action, using local filesystem")
		pipe = pipeline.NewMemoryPipeline()
		fs = afero.NewMemMapFs()

		pr = git.NewMemoryPullRequestProvider([]*git.PullRequest{
			{
				Number: 1,
				Open:   true,
			},
		})
		release = git.NewMemoryReleaseProvider([]*git.Release{})
		repometa = git.NewMemoryRepoMetadataProvider(&git.RemoteRepositoryMetadata{
			Description: "test repo",
			Homepage:    "test.com",
			License:     "Nunya",
		})
	}

	ctx = snake.Bind(ctx, (*git.PullRequestProvider)(nil), pr)
	ctx = snake.Bind(ctx, (*git.ReleaseProvider)(nil), release)
	ctx = snake.Bind(ctx, (*git.RemoteRepositoryMetadataProvider)(nil), repometa)
	ctx = snake.Bind(ctx, (*pipeline.Pipeline)(nil), pipe)
	ctx = snake.Bind(ctx, (*afero.Fs)(nil), fs)
	ctx = snake.Bind(ctx, (*git.GitProvider)(nil), execgit)

	rootCmd := snake.NewRootCommand(ctx, &root.Root{})

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("ERROR: %+v", err)
	}

}
