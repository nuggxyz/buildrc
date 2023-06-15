package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/go-github/v53/github"
	"github.com/rs/zerolog"
)

// all together now
func (me *GithubClient) UploadArtifact(ctx context.Context, file *os.File) (*github.Artifact, error) {
	name := filepath.Base(file.Name())

	art, err := me.CreateWorkflowArtifact(ctx, name)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("artifact", art).Msg("created artifact")

	size, err := me.UploadWorkflowArtifact(ctx, art, file)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Int64("size", int64(size)).Msg("uploaded artifact")

	res, err := me.UpdateWorkflowArtifact(ctx, name, size)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("res", res).Msg("updated artifact")

	return res, nil

}

func (me *GithubClient) UploadWorkflowArtifact(ctx context.Context, artifact string, file *os.File) (int, error) {

	// stat, err := file.Stat()
	// if err != nil {
	// 	return 0, err
	// }

	req, err := http.NewRequest("PUT", artifact+fmt.Sprintf("?itemPath=%s", artifact), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ActionRuntimeToken.Load()))
	req.Header.Set("Accept", "application/json;api-version=6.0-preview")

	chunkSize := 8 * 1024 * 1024 // 8MB
	buffer := make([]byte, chunkSize)
	totalBytesRead := 0
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}

		totalBytesRead += bytesRead

		zerolog.Ctx(ctx).Debug().Int("bytesRead", bytesRead).Msg("uploading artifact")
		req.Header.Set("Content-Length", strconv.FormatInt(int64(bytesRead), 10))
		req.Body = io.NopCloser(bytes.NewBuffer(buffer[:bytesRead]))
		reso, err := http.DefaultClient.Do(req)
		if err != nil {

			zerolog.Ctx(ctx).Error().Err(err).Int("status", reso.StatusCode).Msg("failed to upload artifact")
			return 0, err
		}
		if reso.StatusCode != 200 {
			var interfaceErr interface{}
			err = json.NewDecoder(reso.Body).Decode(&interfaceErr)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Int("status", reso.StatusCode).Msg("failed to upload artifact")
				return 0, err
			}
			zerolog.Ctx(ctx).Error().Any("response_body", interfaceErr).Int("status", reso.StatusCode).Msg("failed to upload artifact")
			return 0, fmt.Errorf("failed to upload artifact: %d", reso.StatusCode)
		}
	}

	return totalBytesRead, nil
}

// Function to create an artifact
func (me *GithubClient) CreateWorkflowArtifact(ctx context.Context, name string) (string, error) {
	id, err := strconv.ParseInt(string(GitHubRunID.Load()), 10, 64)
	if err != nil {
		return "", err
	}

	// Define your artifact details here
	artifactDetails := map[string]interface{}{
		"type": "actions_storage",
		"name": name,
	}

	artifactDetailsBytes, err := json.Marshal(artifactDetails)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/_apis/pipelines/workflows/%d/artifacts?api-version=6.0-preview", ActionRuntimeURL.Load(), id), bytes.NewBuffer(artifactDetailsBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ActionRuntimeToken.Load()))
	req.Header.Set("Accept", "application/json;api-version=6.0-preview")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var artifact struct {
		URL string `json:"fileContainerResourceUrl"`
	}
	err = json.NewDecoder(resp.Body).Decode(&artifact)
	if err != nil {
		return "", err
	}

	return artifact.URL, nil
}

// Function to update an artifact
func (me *GithubClient) UpdateWorkflowArtifact(ctx context.Context, name string, size int) (*github.Artifact, error) {
	id, err := strconv.ParseInt(string(GitHubRunID.Load()), 10, 64)
	if err != nil {
		return nil, err
	}

	// Define your update details here
	updateDetails := map[string]interface{}{
		"size": size,
	}

	updateDetailsBytes, err := json.Marshal(updateDetails)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/_apis/pipelines/workflows/%d/artifacts?api-version=6.0-preview&artifactName=%s", ActionRuntimeURL.Load(), id, name), bytes.NewBuffer(updateDetailsBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ActionRuntimeToken.Load()))
	req.Header.Set("Accept", "application/json;api-version=6.0-preview")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var artifact github.Artifact
	err = json.NewDecoder(resp.Body).Decode(&artifact)
	if err != nil {
		return nil, err
	}

	return &artifact, nil
}
