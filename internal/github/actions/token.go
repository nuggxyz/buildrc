package actions

import "errors"

func (me *GithubActionPipeline) GithubRestApiToken() (string, error) {
	res := EnvVarGithubToken.Load()
	if res == "" {
		return "", errors.New("env var not set or empty: " + string(EnvVarGithubToken))
	}

	return res, nil
}
