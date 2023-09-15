package buildrc_test

import (
	"context"
	"testing"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/walteh/buildrc/gen/mockery"
	"github.com/walteh/buildrc/pkg/buildrc"
)

func TestGetRepo(t *testing.T) {
	tests := []struct {
		name          string
		remoteURL     string
		mockError     error
		expectedOrg   string
		expectedRepo  string
		expectedError error
	}{
		{
			name:          "valid https remote",
			remoteURL:     "https://github.com/org/repo.git",
			expectedOrg:   "org",
			expectedRepo:  "repo",
			expectedError: nil,
		},
		{
			name:          "valid ssh remote",
			remoteURL:     "git@github.com:org/repo.git",
			expectedOrg:   "org",
			expectedRepo:  "repo",
			expectedError: nil,
		},
		{
			name:          "invalid remote",
			remoteURL:     "invalid",
			expectedError: buildrc.ErrCouldNotParseRemoteURL,
		},
		{
			name:          "git provider error",
			mockError:     errors.Errorf("some error"),
			expectedError: errors.Errorf("some error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			mockGitProvider := mockery.NewMockGitProvider_git(t)

			mockGitProvider.EXPECT().GetRemoteURL(mock.Anything).Return(test.remoteURL, test.mockError)

			org, repo, err := buildrc.GetRepo(context.TODO(), mockGitProvider)

			assert.Equal(t, test.expectedOrg, org)
			assert.Equal(t, test.expectedRepo, repo)

			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockGitProvider.AssertExpectations(t)
		})
	}
}
