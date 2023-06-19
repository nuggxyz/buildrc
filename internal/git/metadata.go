package git

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

type RepositoryMetadata struct {
	Owner       string
	Name        string
	Description string
	Homepage    string
	License     string
}

type CommitMetadata struct {
	Branch      string
	Tag         *semver.Version
	Head        string
	ContentHash string
}

type RepositoryMetadataProvider interface {
	GetRepositoryMetadata(ctx context.Context) (*RepositoryMetadata, error)
}

type DockerBakeTemplateTags []string

func BuildDockerBakeTemplateTags(ctx context.Context, repo RepositoryMetadataProvider, comt GitProvider) (DockerBakeTemplateTags, error) {

	commitMetadata, err := GetCommitMetadata(ctx, comt, "HEAD")
	if err != nil {
		return nil, err
	}

	tagnov := strings.TrimPrefix(commitMetadata.Tag.String(), "v")

	strs := []string{}
	strs = append(strs, "type=ref,event=branch")
	strs = append(strs, "type=ref,event=pr")
	strs = append(strs, "type=schedule")
	strs = append(strs, fmt.Sprintf("type=semver,pattern=v{{version}},value=%s", tagnov))
	strs = append(strs, "type=sha")
	strs = append(strs, fmt.Sprintf("type=raw,value=latest,enable=%v", commitMetadata.Branch == "main"))
	strs = append(strs, fmt.Sprintf("type=semver,pattern=v{{major}}.{{minor}},value=%s,enable=%v", tagnov, commitMetadata.Branch == "main"))
	strs = append(strs, fmt.Sprintf("type=semver,pattern=v{{major}},value=%s,enable=%v", tagnov, commitMetadata.Branch == "main"))

	return strs, nil
}

func (me DockerBakeTemplateTags) NewLineString() (string, error) {
	strs := strings.Join([]string(me), "\n")
	res, err := json.Marshal(strs)
	if err != nil {
		return "", err
	}
	return string(res), nil

}

type DockerBakeLabels map[string]string

func BuildDockerBakeLabels(ctx context.Context, name string, repo RepositoryMetadataProvider, comt GitProvider) (DockerBakeLabels, error) {

	commitMetadata, err := GetCommitMetadata(ctx, comt, "HEAD")
	if err != nil {
		return nil, err
	}

	repoMetadata, err := repo.GetRepositoryMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"org.opencontainers.image.title":         name,
		"org.opencontainers.image.source":        repoMetadata.Homepage,
		"org.opencontainers.image.url":           repoMetadata.Homepage,
		"org.opencontainers.image.documentation": repoMetadata.Homepage + "/README.md",
		"org.opencontainers.image.version":       commitMetadata.Tag.String(),
		"org.opencontainers.image.revision":      commitMetadata.Head,
		"org.opencontainers.image.vendor":        repoMetadata.Owner,
		"org.opencontainers.image.licenses":      repoMetadata.License,
		"org.opencontainers.image.created":       time.Now().Format(time.RFC3339),
		"org.opencontainers.image.authors":       repoMetadata.Owner,
		"org.opencontainers.image.ref.name":      commitMetadata.Tag.String(),
		"org.opencontainers.image.description":   repoMetadata.Description,
	}, nil
}

func (me DockerBakeLabels) NewLineString() (string, error) {
	res, err := json.Marshal(me)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
