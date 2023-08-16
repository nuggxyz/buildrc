package actions

import (
	"github.com/nuggxyz/buildrc/internal/env"
)

type EnvVar string

// Github Actions environment variables.
const (
	EnvVarCI                      EnvVar = "CI"
	EnvVarGithubAction            EnvVar = "GITHUB_ACTION"
	EnvVarGithubActionPath        EnvVar = "GITHUB_ACTION_PATH"
	EnvVarGithubActionRepository  EnvVar = "GITHUB_ACTION_REPOSITORY"
	EnvVarGithubActions           EnvVar = "GITHUB_ACTIONS"
	EnvVarGithubActor             EnvVar = "GITHUB_ACTOR"
	EnvVarGithubActorID           EnvVar = "GITHUB_ACTOR_ID"
	EnvVarGithubAPIURL            EnvVar = "GITHUB_API_URL"
	EnvVarGithubBaseRef           EnvVar = "GITHUB_BASE_REF"
	EnvVarGithubEnv               EnvVar = "GITHUB_ENV"
	EnvVarGithubEventName         EnvVar = "GITHUB_EVENT_NAME"
	EnvVarGithubEventPath         EnvVar = "GITHUB_EVENT_PATH"
	EnvVarGithubGraphQLURL        EnvVar = "GITHUB_GRAPHQL_URL"
	EnvVarGithubHeadRef           EnvVar = "GITHUB_HEAD_REF"
	EnvVarGithubJob               EnvVar = "GITHUB_JOB"
	EnvVarGithubOutput            EnvVar = "GITHUB_OUTPUT"
	EnvVarGithubPath              EnvVar = "GITHUB_PATH"
	EnvVarGithubRef               EnvVar = "GITHUB_REF"
	EnvVarGithubRefName           EnvVar = "GITHUB_REF_NAME"
	EnvVarGithubRefProtected      EnvVar = "GITHUB_REF_PROTECTED"
	EnvVarGithubRefType           EnvVar = "GITHUB_REF_TYPE"
	EnvVarGithubRepository        EnvVar = "GITHUB_REPOSITORY"
	EnvVarGithubRepositoryID      EnvVar = "GITHUB_REPOSITORY_ID"
	EnvVarGithubRepositoryOwner   EnvVar = "GITHUB_REPOSITORY_OWNER"
	EnvVarGithubRepositoryOwnerID EnvVar = "GITHUB_REPOSITORY_OWNER_ID"
	EnvVarGithubRetentionDays     EnvVar = "GITHUB_RETENTION_DAYS"
	EnvVarGithubRunAttempt        EnvVar = "GITHUB_RUN_ATTEMPT"
	EnvVarGithubRunID             EnvVar = "GITHUB_RUN_ID"
	EnvVarGithubRunNumber         EnvVar = "GITHUB_RUN_NUMBER"
	EnvVarGithubServerURL         EnvVar = "GITHUB_SERVER_URL"
	EnvVarGithubSHA               EnvVar = "GITHUB_SHA"
	EnvVarGithubStepSummary       EnvVar = "GITHUB_STEP_SUMMARY"
	EnvVarGithubWorkflow          EnvVar = "GITHUB_WORKFLOW"
	EnvVarGithubWorkflowRef       EnvVar = "GITHUB_WORKFLOW_REF"
	EnvVarGithubWorkflowSHA       EnvVar = "GITHUB_WORKFLOW_SHA"
	EnvVarGithubWorkspace         EnvVar = "GITHUB_WORKSPACE"
	EnvVarRunnerArch              EnvVar = "RUNNER_ARCH"
	EnvVarRunnerDebug             EnvVar = "RUNNER_DEBUG"
	EnvVarRunnerName              EnvVar = "RUNNER_NAME"
	EnvVarRunnerOS                EnvVar = "RUNNER_OS"
	EnvVarRunnerTemp              EnvVar = "RUNNER_TEMP"
	EnvVarRunnerToolCache         EnvVar = "RUNNER_TOOL_CACHE"
	EnvVarActionRuntimeURL        EnvVar = "ACTIONS_RUNTIME_URL"
	EnvVarActionRuntimeToken      EnvVar = "ACTIONS_RUNTIME_TOKEN"
	EnvVarGithubToken             EnvVar = "GITHUB_TOKEN"
)

// GetEnv is a helper function that returns the value of the environment variable.
func GetActionEnvVar(key EnvVar) string {
	return env.GetOrEmpty(string(key))
}

func (e EnvVar) Load() string {
	return GetActionEnvVar(e)
}
