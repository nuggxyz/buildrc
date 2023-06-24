package buildrc

import "fmt"

type Github struct {
}

func (me *Github) Repository(pkg *Package, org string, name string) string {
	last := name
	if name != pkg.Name {
		last = fmt.Sprintf("%s/%s", name, pkg.Name)
	}
	return "ghcr.io" + "/" + org + "/" + last
}
