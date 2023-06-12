package docker

import (
	"fmt"
	"strings"
)

func BuildXTagString(repo string, tag string) string {
	ismain := strings.Contains(repo, "main")
	str := ""
	str += "type=ref,event=branch\\n"
	str += fmt.Sprintf("type=semver,pattern=v{{version}},value=%s\\n", tag)
	str += "type=sha\\n"
	str += fmt.Sprintf("type=raw,value=latest,enable=%v\\n", ismain)
	str += fmt.Sprintf("type=semver,pattern=v{{major}}.{{minor}},value=%s,enable=%v\\n", tag, ismain)
	str += fmt.Sprintf("type=semver,pattern=v{{major}},value=%s,enable=%v", tag, ismain)
	return str
}
