package buildrc

type Aws struct {
	IamRole   string `yaml:"iam_role" json:"iam_role"`
	AccountID string `yaml:"account" json:"account"`
	Region    string `yaml:"region" json:"region"`
}

func (me *Aws) Repository(pkg *Package, org string, name string) string {
	return me.AccountID + ".dkr.ecr." + me.Region + ".amazonaws.com" + "/" + org + "/" + name + "/" + pkg.Name
}

func (me *Aws) FullIamRole() string {
	return "arn:aws:iam::" + me.AccountID + ":role/" + me.IamRole
}
