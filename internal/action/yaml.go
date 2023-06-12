package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type GithubActionOutput struct {
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}

type GithubActionStep struct {
	Run   string            `yaml:"run,omitempty"`
	Shell string            `yaml:"shell,omitempty"`
	Id    string            `yaml:"id,omitempty"`
	Name  string            `yaml:"name,omitempty"`
	If    string            `yaml:"if,omitempty"`
	Uses  string            `yaml:"uses,omitempty"`
	With  map[string]string `yaml:"with,omitempty"`
	Env   map[string]string `yaml:"env,omitempty"`
}

type GithubAction struct {
	Name        string                        `yaml:"name"`
	Description string                        `yaml:"description"`
	Inputs      map[string]GithubActionInput  `yaml:"inputs"`
	Outputs     map[string]GithubActionOutput `yaml:"outputs,omitempty"`
	Runs        GithubActionRuns              `yaml:"runs"`
}

type GithubActionRuns struct {
	Using string             `yaml:"using"`
	Image string             `yaml:"image,omitempty"`
	Steps []GithubActionStep `yaml:"steps,omitempty"`
}

type GithubActionInput struct {
	Description string `yaml:"description"`
	Default     string `yaml:"default,omitempty"`
	Required    bool   `yaml:"required"`
}

const INSTALL_SCRIPT = `#!/bin/bash

DOCKER_IMAGE=ghcr.io/nuggxyz/actions/cli:latest
TARGET=${{ inputs.executable_name }}
REPO=${{ github.repository }}

if [ ! -f "$TARGET" ]; then
	if [ "$REPO" == "nuggxyz/actions" ]; then
		echo "inside ci repo, building executable..."
		GOFLAGS="-mod=vendor" go build -v -o ./$TARGET ./cmd/cli/main.go
		exit 0
	fi
	echo "Downloading cli from Docker image..."
	docker pull $DOCKER_IMAGE

	echo "Creating temporary container to extract executable..."
	CONTAINER_ID=$(docker create $DOCKER_IMAGE)
	docker cp $CONTAINER_ID:/main ./$TARGET

	echo "Removing temporary container..."
	docker rm $CONTAINER_ID

	echo "Making the extracted executable executable..."
	chmod +x ./$TARGET
else
  echo "App directory found, skipping download."
fi
`

func (c Command) ToGithubAction() GithubAction {
	inputs := map[string]GithubActionInput{}

	for _, input := range c.Flags {
		if toDefault(input.Name) != input.Name {
			continue
		}
		rep := strings.ReplaceAll(input.Name, "-", "_")

		inputs[rep] = GithubActionInput{
			Description: input.Description,
			Required:    input.Required,
			Default:     input.Default,
		}
	}

	for _, input := range c.Positional {
		if toDefault(input.Name) != input.Name {
			continue
		}
		rep := strings.ReplaceAll(input.Name, "-", "_")
		inputs[rep] = GithubActionInput{
			Description: input.Description,
			Required:    input.Required,
			Default:     input.Default,
		}
	}

	inputs["executable_name"] = GithubActionInput{
		Description: "Target executable name",
		Default:     "nuggxyz_actions_executable",
		Required:    false,
	}

	inputs["action_token"] = GithubActionInput{
		Description: "Github token",
		Required:    true,
	}

	inputs["pkg_read_token"] = GithubActionInput{
		Description: "Github token",
		Required:    true,
	}

	ifcheck := "${{ github.repository == 'nuggxyz/ci' }}"

	return GithubAction{
		Name:        c.Name,
		Description: c.Description,
		Inputs:      inputs,
		Outputs: map[string]GithubActionOutput{
			"result": {
				Description: "Base64 encoded output",
				Value:       "${{ steps.main.outputs.result }}",
			},
		},
		Runs: GithubActionRuns{
			Using: "composite",
			Steps: []GithubActionStep{
				{
					If:   ifcheck,
					Uses: "actions/setup-go@v3",
				},
				{
					If:    ifcheck,
					Run:   "ls -la",
					Shell: "bash",
				},
				{
					Uses: "docker/login-action@v2",
					Id:   "docker-login",
					With: map[string]string{
						"registry": "ghcr.io",
						"username": "${{ github.actor }}",
						"password": "${{ inputs.pkg_read_token }}",
					},
				},
				{
					Run:   INSTALL_SCRIPT,
					Id:    "install",
					Shell: "bash",
					Name:  "Install the executable",
				},
				{
					Run:   c.buildGithubActionRun("./${{ inputs.executable_name }}"),
					Shell: "bash",
					Id:    "main",
					Env: map[string]string{
						"AWS_ACCESS_KEY_ID":     "${{ env.AWS_ACCESS_KEY_ID }}",
						"AWS_SECRET_ACCESS_KEY": "${{ env.AWS_SECRET_ACCESS_KEY }}",
						"AWS_SESSION_TOKEN":     "${{ env.AWS_SESSION_TOKEN }}",
						"AWS_REGION":            "${{ inputs.aws_region }}",
						"AWS_DEFAULT_REGION":    "${{ inputs.aws_region }}",
					},
				},
			},
		},
	}
}

func (action GithubAction) WriteYamlToDir(dir string) error {

	if dir == "" {
		dir = ".github/actions"
	}

	dir = filepath.Join(dir, action.Name)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating actions directory: %v\n", err)
		return err
	}

	actionPath := filepath.Join(dir, "action.yaml")
	actionContent, err := yaml.Marshal(action)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling YAML: %v\n", err)
		return err
	}

	// Write action.yml
	err = os.WriteFile(actionPath, actionContent, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", actionPath, err)
		return err
	}

	return nil
}
