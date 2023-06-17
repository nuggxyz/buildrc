package container

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"
)

const (
	CommandID = "container"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
	Name string `flag:"name" type:"string" default:"main"`
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Invoke(ctx, cp)
	return err
}

func (me *Handler) Invoke(ctx context.Context, cp provider.ContentProvider) (out *any, err error) {
	return me.invoke(ctx, cp)
}

func (me *Handler) invoke(ctx context.Context, r provider.ContentProvider) (out *any, err error) {

	brc, err := load.NewHandler(me.File).Load(ctx, r)
	if err != nil {
		return nil, err
	}

	stup, err := setup.NewHandler("", "").Invoke(ctx, r)
	if err != nil {
		return nil, err
	}

	pkg, ok := brc.PackageByName()[me.Name]
	if !ok {
		return nil, fmt.Errorf("package %s not found", me.Name)
	}

	if len(pkg.DockerPlatforms) == 0 {
		export := map[string]string{
			"BUILDRC_CONTAINER_PUSH": "0",
		}

		err = provider.AddContentToEnv(ctx, r, CommandID, export)
		if err != nil {
			return nil, err
		}

		return

	}

	ghcli, err := github.NewGithubClient(ctx, "", "")
	if err != nil {
		return nil, err
	}

	tags, err := ghcli.BuildXTagString(ctx, stup.Tag)
	if err != nil {
		return nil, err
	}

	labs, err := ghcli.BuildxLabelString(ctx, pkg.Name, stup.Tag)
	if err != nil {
		return nil, err
	}

	img, err := brc.ImagesCSVJSON(pkg, ghcli.OrgName(), ghcli.RepoName())
	if err != nil {
		return nil, err
	}

	cd, err := buildrc.BuildrcCacheDir.Load()
	if err != nil {
		return nil, err
	}

	ccc, err := pkg.DockerBuildArgsCSV()
	if err != nil {
		return nil, err
	}

	uploadToAws := "0"

	if brc.Aws != nil {
		uploadToAws = "1"
	}

	export := map[string]string{
		"BUILDRC_CONTAINER_PUSH":               "1",
		"BUILDRC_CONTAINER_IMAGES_JSON_STRING": img,
		"BUILDRC_CONTAINER_LABELS_JSON_STRING": labs,
		"BUILDRC_CONTAINER_TAGS_JSON_STRING":   tags,
		"BUILDRC_CONTAINER_CONTEXT":            cd,
		"BUILDRC_CONTAINER_DOCKERFILE":         pkg.Dockerfile(),
		"BUILDRC_CONTAINER_PLATFORMS_CSV":      pkg.DockerPlatformsCSV(),
		"BUILDRC_CONTAINER_BUILD_ARGS_CSV":     ccc,
		"BUILDRC_CONTAINER_UPLOAD_TO_AWS":      uploadToAws,
	}

	if brc.Aws != nil {
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_IAM_ROLE"] = brc.Aws.FullIamRole()
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_REGION"] = brc.Aws.Region
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_ACCOUNT_ID"] = brc.Aws.AccountID
		export["BUILDRC_CONTAINER_UPLOAD_TO_AWS_REPOSITORY"] = brc.Aws.Repository(pkg, ghcli.OrgName(), ghcli.RepoName())
	}

	err = provider.AddContentToEnv(ctx, r, CommandID, export)
	if err != nil {
		return nil, err
	}

	return
}
