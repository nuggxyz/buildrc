package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (client *GitHubAtifactClient) CreateAndUploadArtifact(artifactName, content string) error {
	headers := map[string]string{
		"Accept":        "application/json;api-version=6.0-preview",
		"Authorization": "Bearer " + client.RuntimeToken,
	}

	artifactBase := fmt.Sprintf("%s_apis/pipelines/workflows/%s/artifacts?api-version=6.0-preview", client.RuntimeURL, client.RunID)

	resourceURL, err := client.createArtifact(artifactBase, artifactName, headers)
	if err != nil {
		return fmt.Errorf("creating artifact failed: %w", err)
	}

	if err = client.uploadArtifact(resourceURL, artifactName, content, headers); err != nil {
		return fmt.Errorf("uploading artifact failed: %w", err)
	}

	if err = client.updateArtifact(artifactBase, artifactName, len(content), headers); err != nil {
		return fmt.Errorf("updating artifact failed: %w", err)
	}

	return nil
}

func (client *GitHubAtifactClient) createArtifact(url, name string, headers map[string]string) (string, error) {
	postData := map[string]string{
		"type": "actions_storage",
		"name": name,
	}

	jsonValue, _ := json.Marshal(postData)

	resp, err := client.doRequest("POST", url, bytes.NewBuffer(jsonValue), headers)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	return result["fileContainerResourceUrl"].(string), nil
}

func (client *GitHubAtifactClient) uploadArtifact(url, name, content string, headers map[string]string) error {
	headers["Content-Type"] = "application/octet-stream"
	headers["Content-Range"] = fmt.Sprintf("bytes 0-%d/%d", len(content)-1, len(content))

	_, err := client.doRequest("PUT", fmt.Sprintf("%s?itemPath=%s/data.txt", url, name), bytes.NewBufferString(content), headers)
	return err
}

func (client *GitHubAtifactClient) updateArtifact(url, name string, size int, headers map[string]string) error {
	patchData := map[string]int{
		"size": size,
	}

	jsonValue, _ := json.Marshal(patchData)

	_, err := client.doRequest("PATCH", fmt.Sprintf("%s&artifactName=%s", url, name), bytes.NewBuffer(jsonValue), headers)
	return err
}

func (client *GitHubAtifactClient) doRequest(method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
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
