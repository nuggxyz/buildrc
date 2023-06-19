package common

import (
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
)

type Provider interface {
	Git() git.GitProvider
	Release() git.ReleaseProvider
	Pipeline() pipeline.Pipeline
	PR() git.PullRequestProvider
	Buildrc() *buildrc.Buildrc
	RepositoryMetadata() git.RepositoryMetadataProvider
}

type providerGroup struct {
	gi   git.GitProvider
	rel  git.ReleaseProvider
	cp   pipeline.Pipeline
	pr   git.PullRequestProvider
	brc  *buildrc.Buildrc
	meta git.RepositoryMetadataProvider
}

func NewProvider(gi git.GitProvider, rel git.ReleaseProvider, cp pipeline.Pipeline, pr git.PullRequestProvider, brc *buildrc.Buildrc, meta git.RepositoryMetadataProvider) Provider {
	return &providerGroup{gi, rel, cp, pr, brc, meta}
}

func (me *providerGroup) Git() git.GitProvider {
	return me.gi
}

func (me *providerGroup) Release() git.ReleaseProvider {

	return me.rel
}

func (me *providerGroup) Pipeline() pipeline.Pipeline {
	return me.cp
}

func (me *providerGroup) PR() git.PullRequestProvider {
	return me.pr
}

func (me *providerGroup) Buildrc() *buildrc.Buildrc {
	return me.brc
}

func (me *providerGroup) RepositoryMetadata() git.RepositoryMetadataProvider {
	return me.meta
}
