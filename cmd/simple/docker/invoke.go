package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
)

type Handler struct {
	Package string `arg:"package" help:"The name of the package to load."`
	Tag     bool   `flag:"tag" help:"The tag to use for the docker image."`
	Build   bool   `flag:"build" help:"The build to use for the docker image."`
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {

	repom, err := prov.Git().GetLocalRepositoryMetadata(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	pkg, ok := prov.Buildrc().PackageByName()[me.Package]
	if !ok {
		return fmt.Errorf("package %s not found", me.Package)
	}

	should, err := pkg.ShouldBuildDocker(ctx, prov.FileSystem())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	skipDocker := "0"
	if !should || !prov.Pipeline().SupportsDocker() {
		skipDocker = "1"
	}

	targetSemver, err := git.CalculateNextPreReleaseTag(ctx, prov.Buildrc(), prov.Git(), prov.PR())
	if err != nil {
		return err
	}

	tags, err := git.BuildDockerBakeTemplateTags(ctx, prov.Git(), targetSemver)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	bstags, err := git.BuildDockerBakeBuildSpecificTemplateTags(ctx, prov.Git(), targetSemver)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	labs, err := git.BuildDockerBakeLabels(ctx, "", prov.RepositoryMetadata(), prov.Git())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	img, err := prov.Buildrc().ImagesCSVJSON(pkg, repom.Owner, repom.Name)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	// cd, err := pipeline.CacheDir(ctx, prov.Pipeline(), prov.FileSystem())
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
	// 	return err
	// }

	// mycd := filepath.Join(cd, pkg.Name)

	// ccc, err := pkg.DockerBuildArgs(ctx, prov.Pipeline(), prov.FileSystem())
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
	// 	return err
	// }

	// dbajs, err := ccc.JSONString()
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
	// 	return err
	// }

	uploadToAws := "0"

	if prov.Buildrc().Aws != nil {
		uploadToAws = "1"
	}

	alreadyExists := "0"

	b, _, err := git.ReleaseAlreadyExists(ctx, prov.Release(), prov.Git())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	if b {
		zerolog.Ctx(ctx).Info().Str("package", pkg.Name).Msg("package already exists")
		alreadyExists = "1"
	}

	lstr, err := labs.NewLineString()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	bsstr, err := bstags.NewLineString()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}
	tstr, err := tags.NewLineString()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	root, err := prov.Pipeline().RootDir(ctx)
	if err != nil {
		return err
	}

	// for _, dp := range pkg.DockerPlatforms {
	// 	opf, err := dp.OutputFile(pkg)
	// 	if err != nil {
	// 		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
	// 		return err
	// 	}

	// 	_, err = prov.FileSystem().Stat(opf)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	skipTag := "1"
	if me.Tag {
		skipTag = "0"
	}

	skipBuild := "1"
	if me.Build {
		skipBuild = "0"
	}

	export := map[string]string{
		"BUILDRC_SKIP_DOCKER":                               skipDocker,
		"BUILDRC_CONTAINER_IMAGES_JSON_STRING":              img,
		"BUILDRC_CONTAINER_LABELS_JSON_STRING":              lstr,
		"BUILDRC_CONTAINER_BUILD_SPECIFIC_TAGS_JSON_STRING": bsstr,
		"BUILDRC_CONTAINER_TAGS_JSON_STRING":                tstr,
		"BUILDRC_CONTAINER_CONTEXT":                         root,
		"BUILDRC_CONTAINER_DOCKERFILE":                      pkg.Dockerfile(),
		"BUILDRC_CONTAINER_PLATFORMS_CSV":                   pkg.DockerPlatformsCSV(),
		"BUILDRC_CONTAINER_BUILD_ARGS_JSON_STRING":          "",
		"BUILDRC_CONTAINER_UPLOAD_TO_AWS":                   uploadToAws,
		"BUILDRC_CONTAINER_BUILD_EXISTS":                    alreadyExists,
		"BUILDRC_SKIP_DOCKER_BUILD":                         skipBuild,
		"BUILDRC_SKIP_DOCKER_TAG":                           skipTag,
	}

	for a, z := range pkg.UsesMap() {
		export["BUILDRC_USES_"+strings.ToUpper(a)] = z
	}

	if prov.Buildrc().Aws != nil {
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_IAM_ROLE"] = prov.Buildrc().Aws.FullIamRole()
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_REGION"] = prov.Buildrc().Aws.Region
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_ACCOUNT_ID"] = prov.Buildrc().Aws.AccountID
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_REPOSITORY"] = prov.Buildrc().Aws.Repository(pkg, repom.Owner, repom.Name)
	}

	err = pipeline.AddContentToEnvButDontCache(ctx, prov.Pipeline(), prov.FileSystem(), "docker", export)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error is here")
		return err
	}

	return
}
