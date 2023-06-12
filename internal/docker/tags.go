package docker

import (
	"context"
	"fmt"
	"strings"
)

func BuildXTagString(ctx context.Context, repo string, tag string) (string, error) {
	ismain := strings.Contains(repo, "main")
	str := ""
	str += "type=ref,event=branch|"
	str += fmt.Sprintf("type=semver,pattern=v{{version}},value=%s|", tag)
	str += "type=sha|"
	str += fmt.Sprintf("type=raw,value=latest,enable=%v|", ismain)
	str += fmt.Sprintf("type=semver,pattern=v{{major}}.{{minor}},value=%s,enable=%v|", tag, ismain)
	str += fmt.Sprintf("type=semver,pattern=v{{major}},value=%s,enable=%v", tag, ismain)

	return string(str), nil
}
