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
	str += fmt.Sprintf("type=semver,pattern=v{{major}}.{{minor}},value=v%s,enable=%v|", tag, ismain)
	str += fmt.Sprintf("type=semver,pattern=v{{major}},value=v%s,enable=%v", tag, ismain)

	return string(str), nil
}

// type=ref,event=branch
//   type=semver,pattern=v{{version}},value=
//   type=sha
//   type=raw,value=latest,enable=true
//   type=semver,pattern=v{{major}}.{{minor}},value=,enable=true
//   type=semver,pattern=v{{major}},value=,enable=true

//   type=ref,event=branch
//   type=semver,pattern=v{{version}},value=0.1.0
//   type=sha
//   type=raw,value=latest,enable=false
//   type=semver,pattern=v{{major}}.{{minor}},value=0.1.0,enable=false
//   type=semver,pattern=v{{major}},value=0.1.0,enable=false
