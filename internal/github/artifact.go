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

	"golang.org/x/net/context/ctxhttp"
)

type GitHubAtifactClient struct {
	RuntimeToken string
	RuntimeURL   string
	RunID        string
}

func NewGitHubArtifactClientFromEnv() *GitHubAtifactClient {

	return &GitHubAtifactClient{
		RuntimeToken: ActionRuntimeToken.Load(),
		RuntimeURL:   ActionRuntimeURL.Load(),
		RunID:        GitHubRunID.Load(),
	}
}

func (client *GitHubAtifactClient) CreateAndUploadArtifactFile(ctx context.Context, content *os.File) error {
	// Read the entire content of the file
	byteContent, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("reading file content failed: %w", err)
	}

	// Convert the byte slice to a string
	strContent := string(byteContent)

	artifactName := filepath.Base(content.Name())

	// Use the string content in the rest of your method
	// ...

	return client.CreateAndUploadArtifact(ctx, artifactName, strContent)
}

func (client *GitHubAtifactClient) CreateAndUploadArtifact(ctx context.Context, artifactName, content string) error {
	headers := map[string]string{
		"Accept":        "application/json;api-version=6.0-preview",
		"Authorization": "Bearer " + client.RuntimeToken,
	}

	artifactBase := fmt.Sprintf("%s_apis/pipelines/workflows/%s/artifacts?api-version=6.0-preview", client.RuntimeURL, client.RunID)

	resourceURL, err := client.createArtifact(ctx, artifactBase, artifactName, headers)
	if err != nil {
		return fmt.Errorf("creating artifact failed: %w", err)
	}

	if err = client.uploadArtifact(ctx, resourceURL, artifactName, content, headers); err != nil {
		return fmt.Errorf("uploading artifact failed: %w", err)
	}

	if err = client.updateArtifact(ctx, artifactBase, artifactName, len(content), headers); err != nil {
		return fmt.Errorf("updating artifact failed: %w", err)
	}

	return nil
}

func (client *GitHubAtifactClient) createArtifact(ctx context.Context, url, name string, headers map[string]string) (string, error) {
	postData := map[string]string{
		"type": "actions_storage",
		"name": name,
	}

	jsonValue, _ := json.Marshal(postData)

	resp, err := client.doRequest(ctx, "POST", url, bytes.NewBuffer(jsonValue), headers)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	return result["fileContainerResourceUrl"].(string), nil
}

func (client *GitHubAtifactClient) uploadArtifact(ctx context.Context, url, name, content string, headers map[string]string) error {
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Range"] = fmt.Sprintf("bytes 0-%d/%d", len(content)-1, len(content))

	_, err := client.doRequest(ctx, "PUT", fmt.Sprintf("%s?itemPath=%s/data.txt", url, name), bytes.NewBufferString(content), headers)
	return err
}

func (client *GitHubAtifactClient) updateArtifact(ctx context.Context, url, name string, size int, headers map[string]string) error {
	patchData := map[string]int{
		"size": size,
	}

	jsonValue, _ := json.Marshal(patchData)

	_, err := client.doRequest(ctx, "PATCH", fmt.Sprintf("%s&artifactName=%s", url, name), bytes.NewBuffer(jsonValue), headers)
	return err
}

func (client *GitHubAtifactClient) doRequest(ctx context.Context, method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := ctxhttp.Do(ctx, http.DefaultClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Usage:
// GitHubAtifactClient := &GitHubAtifactClient{
//	RuntimeToken: "YOUR_RUNTIME_TOKEN",
//	RuntimeURL:   "YOUR_RUNTIME_URL",
//	RunID:        "YOUR_RUN_ID",
//}
// err := githubClient.CreateAndUploadArtifact("artifact-name", "hello world")
// if err != nil {
//	fmt.Println(err)
//}
