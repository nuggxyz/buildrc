package restapi

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/spf13/afero"
)

var _ git.ReleaseProvider = (*GithubClient)(nil)

func (me *GithubClient) CreateRelease(ctx context.Context, g git.GitProvider, t *semver.Version) (*git.Release, error) {

	cmt, err := g.GetCurrentCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	rel, _, err := me.Client().Repositories.CreateRelease(ctx, me.OrgName(), me.RepoName(), &github.RepositoryRelease{
		TargetCommitish: &cmt,
		Name:            github.String(t.String() + " draft"),
		TagName:         github.String(t.String()),

		Draft:      github.Bool(true),
		Prerelease: github.Bool(t.Prerelease() != ""),
	})

	if err != nil {
		return nil, err
	}

	return &git.Release{
		ID:         fmt.Sprintf("%d", rel.GetID()),
		CommitHash: cmt,
		Tag:        getUntaggedTagFromRelease(rel),
		Artifacts:  []string{},
	}, nil

}

func (me *GithubClient) UploadReleaseArtifact(ctx context.Context, r *git.Release, name string, file afero.File) error {

	inter, err := strconv.Atoi(r.ID)
	if err != nil {
		return err
	}

	upload, _, err := me.UploadReleaseAsset(ctx, me.OrgName(), me.RepoName(), int64(inter), name, file)
	if err != nil {
		return err
	}

	r.Artifacts = append(r.Artifacts, upload.GetName())

	return nil

}

func (s *GithubClient) UploadReleaseAsset(ctx context.Context, owner, repo string, id int64, str string, file afero.File) (*github.ReleaseAsset, *github.Response, error) {
	u := fmt.Sprintf("repos/%s/%s/releases/%d/assets?name=%s", owner, repo, id, str)

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	if stat.IsDir() {
		return nil, nil, errors.New("the asset to upload can't be a directory")
	}

	mediaType := mime.TypeByExtension(filepath.Ext(file.Name()))

	req, err := s.client.NewUploadRequest(u, file, stat.Size(), mediaType)
	if err != nil {
		return nil, nil, err
	}

	asset := new(github.ReleaseAsset)
	resp, err := s.client.Do(ctx, req, asset)
	if err != nil {
		return nil, resp, err
	}
	return asset, resp, nil
}

func (me *GithubClient) DownloadReleaseArtifact(ctx context.Context, r *git.Release, name string, fs afero.Fs) (afero.File, error) {

	inter, err := strconv.Atoi(r.ID)
	if err != nil {
		return nil, err
	}

	assets, _, err := me.Client().Repositories.ListReleaseAssets(ctx, me.OrgName(), me.RepoName(), int64(inter), &github.ListOptions{
		PerPage: 100,
	})

	if err != nil {
		return nil, err
	}

	var id int64

	for _, v := range assets {
		if v.GetName() == name {
			id = v.GetID()
		}
	}

	if id == 0 {
		return nil, errors.New("no asset found")
	}

	res, _, err := me.Client().Repositories.DownloadReleaseAsset(ctx, me.OrgName(), me.RepoName(), id, nil)
	if err != nil {
		return nil, err
	}

	fle, err := afero.TempFile(fs, "", name)
	if err != nil {

		return nil, err
	}

	err = afero.WriteReader(fs, fle.Name(), res)
	if err != nil {
		return nil, err
	}

	return fle, nil
}

func (me *GithubClient) GetReleaseByCommit(ctx context.Context, ref string) (*git.Release, error) {
	return me.GetReleaseByTag(ctx, ref)
}

func (me *GithubClient) TagRelease(ctx context.Context, r *git.Release, vers *semver.Version, commit string) (*git.Release, error) {

	inter, err := strconv.Atoi(r.ID)
	if err != nil {
		return nil, err
	}

	rel, _, err := me.Client().Repositories.EditRelease(ctx, me.OrgName(), me.RepoName(), int64(inter), &github.RepositoryRelease{
		TagName:         github.String(vers.String()),
		TargetCommitish: github.String(commit),
		Draft:           github.Bool(false),
		Prerelease:      github.Bool(vers.Prerelease() != ""),
	})

	if err != nil {
		return nil, err
	}

	return &git.Release{
		ID:         fmt.Sprintf("%d", rel.GetID()),
		CommitHash: rel.GetTargetCommitish(),
		Tag:        rel.GetTagName(),
		Artifacts:  []string{},
		Draft:      rel.GetDraft(),
	}, nil

}

func getUntaggedTagFromRelease(rel *github.RepositoryRelease) string {
	arr := strings.Split(rel.GetHTMLURL(), "/")
	return arr[len(arr)-1]
}

func (me *GithubClient) GetReleaseByTag(ctx context.Context, tag string) (*git.Release, error) {

	rel, _, err := me.Client().Repositories.GetReleaseByTag(ctx, me.OrgName(), me.RepoName(), tag)
	if err != nil {
		return nil, err
	}

	id := rel.GetTagName()

	if id == "" {
		id = getUntaggedTagFromRelease(rel)
	}

	return &git.Release{
		ID:         fmt.Sprintf("%d", rel.GetID()),
		CommitHash: rel.GetTargetCommitish(),
		Tag:        id,
		Artifacts:  []string{},
		Draft:      rel.GetDraft(),
	}, nil
}

func (me *GithubClient) ListRecentReleases(ctx context.Context, limit int) ([]*git.Release, error) {

	rels, _, err := me.Client().Repositories.ListReleases(ctx, me.OrgName(), me.RepoName(), &github.ListOptions{
		PerPage: limit,
	})
	if err != nil {
		return nil, err
	}

	var releases []*git.Release

	for _, v := range rels {
		releases = append(releases, &git.Release{
			ID:         fmt.Sprintf("%d", v.GetID()),
			CommitHash: v.GetTargetCommitish(),
			Tag:        v.GetTagName(),
			Artifacts:  []string{},
			Draft:      v.GetDraft(),
		})
	}

	return releases, nil

}
