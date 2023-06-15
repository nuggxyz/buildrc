package github

import (
	"github.com/nuggxyz/buildrc/internal/env"
)

type GitHubActionEnvVar string

// GitHub Actions environment variables.
const (
	CI                      GitHubActionEnvVar = "CI"
	GitHubAction            GitHubActionEnvVar = "GITHUB_ACTION"
	GitHubActionPath        GitHubActionEnvVar = "GITHUB_ACTION_PATH"
	GitHubActionRepository  GitHubActionEnvVar = "GITHUB_ACTION_REPOSITORY"
	GitHubActions           GitHubActionEnvVar = "GITHUB_ACTIONS"
	GitHubActor             GitHubActionEnvVar = "GITHUB_ACTOR"
	GitHubActorID           GitHubActionEnvVar = "GITHUB_ACTOR_ID"
	GitHubAPIURL            GitHubActionEnvVar = "GITHUB_API_URL"
	GitHubBaseRef           GitHubActionEnvVar = "GITHUB_BASE_REF"
	GitHubEnv               GitHubActionEnvVar = "GITHUB_ENV"
	GitHubEventName         GitHubActionEnvVar = "GITHUB_EVENT_NAME"
	GitHubEventPath         GitHubActionEnvVar = "GITHUB_EVENT_PATH"
	GitHubGraphQLURL        GitHubActionEnvVar = "GITHUB_GRAPHQL_URL"
	GitHubHeadRef           GitHubActionEnvVar = "GITHUB_HEAD_REF"
	GitHubJob               GitHubActionEnvVar = "GITHUB_JOB"
	GitHubPath              GitHubActionEnvVar = "GITHUB_PATH"
	GitHubRef               GitHubActionEnvVar = "GITHUB_REF"
	GitHubRefName           GitHubActionEnvVar = "GITHUB_REF_NAME"
	GitHubRefProtected      GitHubActionEnvVar = "GITHUB_REF_PROTECTED"
	GitHubRefType           GitHubActionEnvVar = "GITHUB_REF_TYPE"
	GitHubRepository        GitHubActionEnvVar = "GITHUB_REPOSITORY"
	GitHubRepositoryID      GitHubActionEnvVar = "GITHUB_REPOSITORY_ID"
	GitHubRepositoryOwner   GitHubActionEnvVar = "GITHUB_REPOSITORY_OWNER"
	GitHubRepositoryOwnerID GitHubActionEnvVar = "GITHUB_REPOSITORY_OWNER_ID"
	GitHubRetentionDays     GitHubActionEnvVar = "GITHUB_RETENTION_DAYS"
	GitHubRunAttempt        GitHubActionEnvVar = "GITHUB_RUN_ATTEMPT"
	GitHubRunID             GitHubActionEnvVar = "GITHUB_RUN_ID"
	GitHubRunNumber         GitHubActionEnvVar = "GITHUB_RUN_NUMBER"
	GitHubServerURL         GitHubActionEnvVar = "GITHUB_SERVER_URL"
	GitHubSHA               GitHubActionEnvVar = "GITHUB_SHA"
	GitHubStepSummary       GitHubActionEnvVar = "GITHUB_STEP_SUMMARY"
	GitHubWorkflow          GitHubActionEnvVar = "GITHUB_WORKFLOW"
	GitHubWorkflowRef       GitHubActionEnvVar = "GITHUB_WORKFLOW_REF"
	GitHubWorkflowSHA       GitHubActionEnvVar = "GITHUB_WORKFLOW_SHA"
	GitHubWorkspace         GitHubActionEnvVar = "GITHUB_WORKSPACE"
	RunnerArch              GitHubActionEnvVar = "RUNNER_ARCH"
	RunnerDebug             GitHubActionEnvVar = "RUNNER_DEBUG"
	RunnerName              GitHubActionEnvVar = "RUNNER_NAME"
	RunnerOS                GitHubActionEnvVar = "RUNNER_OS"
	RunnerTemp              GitHubActionEnvVar = "RUNNER_TEMP"
	RunnerToolCache         GitHubActionEnvVar = "RUNNER_TOOL_CACHE"
	ActionRuntimeURL        GitHubActionEnvVar = "ACTIONS_RUNTIME_URL"
	ActionRuntimeToken      GitHubActionEnvVar = "ACTIONS_RUNTIME_TOKEN"
)

// GetEnv is a helper function that returns the value of the environment variable.
func GetActionEnvVar(key GitHubActionEnvVar) string {
	return env.GetOrEmpty(string(key))
}

func (e GitHubActionEnvVar) Load() string {
	return GetActionEnvVar(e)
}
