package action

import (
	"fmt"
	"strings"
)

type Input struct {
	Description string `yaml:"description,flow"`
	Name        string `yaml:"name,flow"`
	Required    bool   `yaml:"required,flow"`
	Default     string `yaml:"default,omitempty,flow"`
}

type Command struct {
	Name        string  `yaml:"name,flow"`
	Description string  `yaml:"description,flow"`
	Flags       []Input `yaml:"inputs,flow"`
	Positional  []Input `yaml:"positional,flow"`
}

func toDefault(name string) string {
	if name == "access-token" {
		return "action_token"
	}
	return name
}

func (me *Command) buildGithubActionRun(executable string) string {

	cmd := fmt.Sprintf("%s %s", executable, me.Name)

	for _, p := range me.Positional {
		rep := strings.ReplaceAll(toDefault(p.Name), "-", "_")

		cmd += fmt.Sprintf(" ${{ inputs.%s }}", rep)

	}

	for _, f := range me.Flags {
		rep := strings.ReplaceAll(toDefault(f.Name), "-", "_")

		cmd += fmt.Sprintf(" --%s ${{ inputs.%s }}", f.Name, rep)
	}

	cmd += ""

	return cmd

}
