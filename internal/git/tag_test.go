package git_test

// func mustParseSemver(t *testing.T, s string) *semver.Version {
// 	v, err := semver.NewVersion(s)
// 	require.NoError(t, err)
// 	return v
// }

// func TestCalculateNextPreReleaseTag(t *testing.T) {
// 	ctx := context.Background()
// 	brc := &buildrc.Buildrc{Version: 1}

// 	t.Run("GetLatestSemverTagFromRef HEAD fails", func(t *testing.T) {
// 		gitp := &gitmocks.GitProvider{}
// 		prp := &gitmocks.PullRequestProvider{}
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "HEAD").Return(nil, errors.New("error"))
// 		_, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
// 		require.Error(t, err)
// 		mock.AssertExpectationsForObjects(t, gitp, prp)

// 	})

// 	t.Run("GetLatestSemverTagFromRef main fails", func(t *testing.T) {
// 		gitp := &gitmocks.GitProvider{}
// 		prp := &gitmocks.PullRequestProvider{}
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "HEAD").Return(mustParseSemver(t, "1.0.0"), nil)
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "main").Return(nil, errors.New("error"))
// 		_, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
// 		require.Error(t, err)
// 		mock.AssertExpectationsForObjects(t, gitp, prp)

// 	})

// 	t.Run("getLatestPullRequest fails", func(t *testing.T) {
// 		gitp := &gitmocks.GitProvider{}
// 		prp := &gitmocks.PullRequestProvider{}
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "HEAD").Return(mustParseSemver(t, "1.0.0"), nil)
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "main").Return(mustParseSemver(t, "1.0.0"), nil)
// 		prp.EXPECT().ListRecentPullRequests(ctx, "HEAD").Return(nil, errors.New("error"))
// 		_, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
// 		require.Error(t, err)
// 		mock.AssertExpectationsForObjects(t, gitp, prp)

// 	})

// 	t.Run("squash merge commit pr", func(t *testing.T) {
// 		gitp := &gitmocks.GitProvider{}
// 		prp := &gitmocks.PullRequestProvider{}
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "HEAD").Return(mustParseSemver(t, "1.0.0"), nil)
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "main").Return(mustParseSemver(t, "1.0.0"), nil)
// 		gitp.EXPECT().GetCurrentBranchFromRef(ctx, "HEAD").Return("main", nil)
// 		gitp.EXPECT().GetContentHashFromRef(ctx, "HEAD").Return("abc", nil)
// 		gitp.EXPECT().GetContentHashFromRef(ctx, "xyz").Return("abc", nil)
// 		prp.EXPECT().ListRecentPullRequests(ctx, "HEAD").Return([]*git.PullRequest{}, nil)
// 		prp.EXPECT().ListRecentPullRequests(ctx, "main").Return([]*git.PullRequest{
// 			{
// 				Head:   "xyz",
// 				Number: 44,
// 				Closed: true,
// 				Open:   false,
// 			},
// 		}, nil)

// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "xyz").Return(mustParseSemver(t, "2.0.0-prerelease"), nil)

// 		res, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
// 		require.NoError(t, err)
// 		assert.Equal(t, "2.0.0", res.String()) // assuming versioning works this way

// 		mock.AssertExpectationsForObjects(t, gitp, prp)
// 	})

// 	t.Run("add commit to pr", func(t *testing.T) {
// 		gitp := &gitmocks.GitProvider{}
// 		prp := &gitmocks.PullRequestProvider{}
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "HEAD").Return(mustParseSemver(t, "2.0.0-pr.99"), nil)
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "main").Return(mustParseSemver(t, "1.0.0"), nil)
// 		prp.EXPECT().ListRecentPullRequests(ctx, "HEAD").Return([]*git.PullRequest{
// 			{
// 				Head:   "notmaincommit",
// 				Number: 99,
// 				Closed: false,
// 				Open:   true,
// 			},
// 		}, nil)

// 		res, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
// 		require.NoError(t, err)
// 		assert.Equal(t, "2.0.0-pr.99", res.String()) // assuming versioning works this way

// 		mock.AssertExpectationsForObjects(t, gitp, prp)
// 	})

// 	t.Run("create new pr", func(t *testing.T) {
// 		gitp := &gitmocks.GitProvider{}
// 		prp := &gitmocks.PullRequestProvider{}
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "HEAD").Return(mustParseSemver(t, "1.0.0"), nil)
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "main").Return(mustParseSemver(t, "1.0.0"), nil)
// 		prp.EXPECT().ListRecentPullRequests(ctx, "HEAD").Return([]*git.PullRequest{
// 			{
// 				Head:   "notmaincommit",
// 				Number: 99,
// 				Closed: false,
// 				Open:   true,
// 			},
// 		}, nil)

// 		res, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
// 		require.NoError(t, err)
// 		assert.Equal(t, "1.1.0-pr.99", res.String()) // assuming versioning works this way

// 		mock.AssertExpectationsForObjects(t, gitp, prp)
// 	})

// 	t.Run("commit to main", func(t *testing.T) {
// 		gitp := &gitmocks.GitProvider{}
// 		prp := &gitmocks.PullRequestProvider{}
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "HEAD").Return(mustParseSemver(t, "1.0.0"), nil)
// 		gitp.EXPECT().GetLatestSemverTagFromRef(ctx, "main").Return(mustParseSemver(t, "1.0.0"), nil)
// 		gitp.EXPECT().GetCurrentBranchFromRef(ctx, "HEAD").Return("main", nil)
// 		prp.EXPECT().ListRecentPullRequests(ctx, "main").Return([]*git.PullRequest{}, nil)

// 		prp.EXPECT().ListRecentPullRequests(ctx, "HEAD").Return([]*git.PullRequest{}, nil)

// 		gitp.EXPECT().GetContentHashFromRef(ctx, "HEAD").Return("abc", nil)
// 		res, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
// 		require.NoError(t, err)
// 		assert.Equal(t, "1.0.1", res.String()) // assuming versioning works this way

// 		mock.AssertExpectationsForObjects(t, gitp, prp)
// 	})

// }
