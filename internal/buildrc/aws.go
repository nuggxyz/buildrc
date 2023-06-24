package buildrc

import "fmt"

type Aws struct {
	IamRole   string `yaml:"iam_role" json:"iam_role"`
	AccountID string `yaml:"account" json:"account"`
	Region    string `yaml:"region" json:"region"`
}

func (me *Aws) Repository(pkg *Package, org string, name string) string {
	last := name
	if name != pkg.Name {
		last = fmt.Sprintf("%s/%s", name, pkg.Name)
	}
	return me.AccountID + ".dkr.ecr." + me.Region + ".amazonaws.com" + "/" + org + "/" + last
}

func (me *Aws) FullIamRole() string {
	return "arn:aws:iam::" + me.AccountID + ":role/" + me.IamRole
}
