package common

import (
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/spf13/afero"
)

type Provider interface {
	Git() git.GitProvider
	Release() git.ReleaseProvider
	Pipeline() pipeline.Pipeline
	PR() git.PullRequestProvider
	Buildrc() *buildrc.Buildrc
	RepositoryMetadata() git.RemoteRepositoryMetadataProvider
	FileSystem() afero.Fs
}

type providerGroup struct {
	gi   git.GitProvider
	rel  git.ReleaseProvider
	cp   pipeline.Pipeline
	pr   git.PullRequestProvider
	brc  *buildrc.Buildrc
	meta git.RemoteRepositoryMetadataProvider
	fs   afero.Fs
}

func NewProvider(gi git.GitProvider, rel git.ReleaseProvider, cp pipeline.Pipeline, pr git.PullRequestProvider, brc *buildrc.Buildrc, meta git.RemoteRepositoryMetadataProvider, fs afero.Fs) Provider {
	return &providerGroup{gi, rel, cp, pr, brc, meta, fs}
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

func (me *providerGroup) RepositoryMetadata() git.RemoteRepositoryMetadataProvider {
	return me.meta
}

func (me *providerGroup) FileSystem() afero.Fs {
	return me.fs
}
