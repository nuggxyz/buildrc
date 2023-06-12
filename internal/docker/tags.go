package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
)

func BuildXTagString(ctx context.Context, repo string, tag string) (string, error) {
	ismain := strings.Contains(repo, "main")
	str := ""
	str += "type=ref,event=branch\n"
	str += fmt.Sprintf("type=semver,pattern=v{{version}},value=%s\n", tag)
	str += "type=sha\n"
	str += fmt.Sprintf("type=raw,value=latest,enable=%v\n", ismain)
	str += fmt.Sprintf("type=semver,pattern=v{{major}}.{{minor}},value=%s,enable=%v\n", tag, ismain)
	str += fmt.Sprintf("type=semver,pattern=v{{major}},value=%s,enable=%v", tag, ismain)

	cmd := exec.Command("docker", "buildx", "bake", "-set", fmt.Sprintf("*.tags=%s", str))
	output, err := cmd.CombinedOutput()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msgf("Error running docker buildx bake: %s", string(output))
		return "", err
	}
	return string(output), err
}
