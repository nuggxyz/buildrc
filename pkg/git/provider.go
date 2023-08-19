package git

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type GitProvider interface {
	LocalRepositoryMetadataProvider
	GetCurrentShortHashFromRef(ctx context.Context, ref string) (string, error)
	GetCurrentCommitFromRef(ctx context.Context, ref string) (string, error)
	GetCurrentCommitMessageFromRef(ctx context.Context, ref string) (string, error)
	GetCurrentBranchFromRef(ctx context.Context, ref string) (string, error)
	GetLatestSemverTagFromRef(ctx context.Context, ref string) (*semver.Version, error)
	GetContentHashFromRef(ctx context.Context, ref string) (string, error)
}
