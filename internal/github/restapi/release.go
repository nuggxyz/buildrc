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
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

var _ git.ReleaseProvider = (*GithubClient)(nil)

func (me *GithubClient) GetReleaseByID(ctx context.Context, id string) (*git.Release, error) {
	inter, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	rel, _, err := me.Client().Repositories.GetRelease(ctx, me.OrgName(), me.RepoName(), inter)
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

func (me *GithubClient) CreateRelease(ctx context.Context, g git.GitProvider, t *semver.Version) (*git.Release, error) {

	cmt, err := g.GetCurrentCommitFromRef(ctx, "HEAD")
	if err != nil {
		return nil, err
	}

	tag := "v" + t.String()

	rela, _, err := me.Client().Repositories.GetReleaseByTag(ctx, me.OrgName(), me.RepoName(), tag)
	if err == nil {
		// Release already exists
		zerolog.Ctx(ctx).Info().Msgf("Release %s already exists", tag)
		return &git.Release{
			ID:         fmt.Sprintf("%d", rela.GetID()),
			CommitHash: rela.GetTargetCommitish(),
			Tag:        rela.GetTagName(),
			Artifacts:  []string{},
			Draft:      rela.GetDraft(),
		}, nil
	}

	rel, _, err := me.Client().Repositories.CreateRelease(ctx, me.OrgName(), me.RepoName(), &github.RepositoryRelease{
		TargetCommitish: &cmt,
		Name:            github.String(t.String() + " draft"),
		TagName:         github.String(tag),
		Draft:           github.Bool(true),
		Prerelease:      github.Bool(t.Prerelease() != ""),
	})

	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Info().Msgf("Created release %s", tag)

	return &git.Release{
		ID:         fmt.Sprintf("%d", rel.GetID()),
		CommitHash: cmt,
		Tag:        tag,
		Artifacts:  []string{},
		Draft:      rel.GetDraft(),
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

	// seek to the beginning of the file
	_, err := file.Seek(0, 0)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to seek to the beginning of the file")
		return nil, nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	if stat.IsDir() {
		return nil, nil, errors.New("the asset to upload can't be a directory")
	}

	if stat.Size() == 0 {

		zerolog.Ctx(ctx).Warn().Interface("stat", stat).Msg("the asset to upload is empty")
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

func (me *GithubClient) TagRelease(ctx context.Context, prov git.GitProvider, vers *semver.Version) (*git.Release, error) {

	// rels, _, err := me.Client().Repositories.ListReleases(ctx, me.OrgName(), me.RepoName(), &github.ListOptions{
	// 	PerPage: 1000,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// for _, v := range rels {

	// 	isTrash := strings.Contains(v.GetTagName(), vers.String()) || (v.CreatedAt.Before(time.Now().Add(-time.Hour*1)) && v.GetDraft())

	// 	if isTrash {
	// 		zerolog.Ctx(ctx).Info().Msgf("deleting release %s", v.GetTagName())
	// 		_, err = me.Client().Repositories.DeleteRelease(ctx, me.OrgName(), me.RepoName(), v.GetID())
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	if (v.GetTagName() == vers.String() || v.GetTagName() == "v"+vers.String()) && !v.GetDraft() { // if the release is a draft, it can have a tag, but the tag is not applied to the repo
	// 		zerolog.Ctx(ctx).Info().Msgf("deleting tag %s", v.GetTagName())

	// 		_, err = me.Client().Git.DeleteRef(ctx, me.OrgName(), me.RepoName(), fmt.Sprintf("tags/%s", v.GetTagName()))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	tag := fmt.Sprintf("v%s", vers.String())

	cmt, err := prov.GetCurrentCommitFromRef(ctx, "HEAD")
	if err != nil {
		return nil, err
	}

	rel, _, err := me.Client().Repositories.CreateRelease(ctx, me.OrgName(), me.RepoName(), &github.RepositoryRelease{
		TargetCommitish: &cmt,
		Name:            github.String(vers.String()),
		TagName:         github.String(tag),
		// we want prereleases to be drafts so that they manually have to be published
		Draft:      github.Bool(true),
		Prerelease: github.Bool(vers.Prerelease() != ""),
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

func (me *GithubClient) HasReleaseArtifact(ctx context.Context, r *git.Release, name string) (bool, error) {

	inter, err := strconv.Atoi(r.ID)
	if err != nil {
		return false, err
	}

	assets, _, err := me.Client().Repositories.ListReleaseAssets(ctx, me.OrgName(), me.RepoName(), int64(inter), &github.ListOptions{
		PerPage: 100,
	})

	if err != nil {
		return false, err
	}

	for _, v := range assets {
		if v.GetName() == name {
			return true, nil
		}
	}

	return false, nil
}

func (me *GithubClient) DeleteReleaseArtifact(ctx context.Context, r *git.Release, name string) error {

	inter, err := strconv.Atoi(r.ID)
	if err != nil {
		return err
	}

	assets, _, err := me.Client().Repositories.ListReleaseAssets(ctx, me.OrgName(), me.RepoName(), int64(inter), &github.ListOptions{
		PerPage: 100,
	})

	if err != nil {
		return err
	}

	var id int64

	for _, v := range assets {
		if v.GetName() == name {
			id = v.GetID()
		}
	}

	if id == 0 {
		return errors.New("no asset found")
	}

	_, err = me.Client().Repositories.DeleteReleaseAsset(ctx, me.OrgName(), me.RepoName(), id)
	if err != nil {
		return err
	}

	return nil
}

func (me *GithubClient) TakeReleaseOutOfDraft(ctx context.Context, rel *git.Release) error {

	inter, err := strconv.Atoi(rel.ID)
	if err != nil {
		return err
	}

	_, _, err = me.Client().Repositories.EditRelease(ctx, me.OrgName(), me.RepoName(), int64(inter), &github.RepositoryRelease{
		Draft: github.Bool(false),
	})
	if err != nil {
		return err
	}

	return nil
}
