package git

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type GitProvider interface {
	GetCurrentCommitHash(ctx context.Context) (string, error)
	GetCurrentBranch(ctx context.Context) (string, error)
	GetLatestSemverTagFromRef(ctx context.Context, ref string) (*semver.Version, error)
	GetContentHash(ctx context.Context, sha string) (string, error)
}
