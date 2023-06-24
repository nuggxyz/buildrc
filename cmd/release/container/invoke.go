package container

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
)

const (
	CommandID = "container"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
	Name string `arg:"name" help:"The name of the package to load."`
}

func (me *Handler) Run(ctx context.Context, cp common.Provider) (err error) {
	_, err = me.Invoke(ctx, cp)
	return err
}

func (me *Handler) Invoke(ctx context.Context, cp common.Provider) (out *any, err error) {
	return me.invoke(ctx, cp)
}

func (me *Handler) invoke(ctx context.Context, prov common.Provider) (out *any, err error) {

	repom, err := prov.Git().GetLocalRepositoryMetadata(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	pkg, ok := prov.Buildrc().PackageByName()[me.Name]
	if !ok {
		return nil, fmt.Errorf("package %s not found", me.Name)
	}

	if len(pkg.DockerPlatforms) == 0 {
		export := map[string]string{
			"BUILDRC_CONTAINER_PUSH": "0",
		}

		err = pipeline.AddContentToEnvButDontCache(ctx, prov.Pipeline(), prov.FileSystem(), CommandID, export)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
			return nil, err
		}

		return
	}

	ss, err := setup.NewHandler("", "").Invoke(ctx, prov)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	tags, err := git.BuildDockerBakeTemplateTags(ctx, prov.Git(), ss.TagSemver)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	labs, err := git.BuildDockerBakeLabels(ctx, "", prov.RepositoryMetadata(), prov.Git())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	img, err := prov.Buildrc().ImagesCSVJSON(pkg, repom.Owner, repom.Name)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	cd, err := pipeline.CacheDir(ctx, prov.Pipeline(), prov.FileSystem())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	ccc, err := pkg.DockerBuildArgs(ctx, prov.Pipeline(), prov.FileSystem())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	dbajs, err := ccc.JSONString()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	uploadToAws := "0"

	if prov.Buildrc().Aws != nil {
		uploadToAws = "1"
	}

	alreadyExists := "0"

	b, _, err := git.ReleaseAlreadyExists(ctx, prov.Release(), prov.Git())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	if b {
		zerolog.Ctx(ctx).Info().Str("package", pkg.Name).Msg("package already exists")
		alreadyExists = "1"
	}

	lstr, err := labs.NewLineString()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	tstr, err := tags.NewLineString()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	export := map[string]string{
		"BUILDRC_CONTAINER_PUSH":                   "1",
		"BUILDRC_CONTAINER_IMAGES_JSON_STRING":     img,
		"BUILDRC_CONTAINER_LABELS_JSON_STRING":     lstr,
		"BUILDRC_CONTAINER_TAGS_JSON_STRING":       tstr,
		"BUILDRC_CONTAINER_CONTEXT":                cd,
		"BUILDRC_CONTAINER_DOCKERFILE":             pkg.Dockerfile(),
		"BUILDRC_CONTAINER_PLATFORMS_CSV":          pkg.DockerPlatformsCSV(),
		"BUILDRC_CONTAINER_BUILD_ARGS_JSON_STRING": dbajs,
		"BUILDRC_CONTAINER_UPLOAD_TO_AWS":          uploadToAws,
		"BUILDRC_CONTAINER_BUILD_EXISTS":           alreadyExists,
	}

	if prov.Buildrc().Aws != nil {
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_IAM_ROLE"] = prov.Buildrc().Aws.FullIamRole()
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_REGION"] = prov.Buildrc().Aws.Region
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_ACCOUNT_ID"] = prov.Buildrc().Aws.AccountID
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_REPOSITORY"] = prov.Buildrc().Aws.Repository(pkg, repom.Owner, repom.Name)
	}

	err = pipeline.AddContentToEnvButDontCache(ctx, prov.Pipeline(), prov.FileSystem(), CommandID, export)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return nil, err
	}

	return
}
