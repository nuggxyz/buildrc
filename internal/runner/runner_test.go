package runner_test

import (
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/file"
	"github.com/nuggxyz/buildrc/internal/logging"
	"github.com/nuggxyz/buildrc/internal/runner"

	"context"
	"errors"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func mustMarshalYAML(obj interface{}) []byte {
	b, err := yaml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}

func TestGHActionContentProvider(t *testing.T) {

	testCases := []struct {
		name             string
		envVars          map[string]string
		expectedErr      error
		cmdID            string
		saveData         []byte
		expectedLoadData []byte
	}{
		{
			name: "Valid GitHub Action",
			envVars: map[string]string{
				"CI":             "true",
				"GITHUB_ACTIONS": "true",
				"GITHUB_ENV":     "test_env",
				"GITHUB_OUTPUT":  "test_output",
				"RUNNER_TEMP":    "test_temp",
			},
			expectedErr:      nil,
			cmdID:            "123",
			saveData:         []byte("test save data"),
			expectedLoadData: []byte("test save data"),
		},
		{
			name: "Valid GitHub Action with JSON",
			envVars: map[string]string{
				"CI":             "true",
				"GITHUB_ACTIONS": "true",
				"GITHUB_ENV":     "test_env",
				"GITHUB_OUTPUT":  "test_output",
				"RUNNER_TEMP":    "test_temp",
			},
			expectedErr:      nil,
			cmdID:            "123",
			saveData:         []byte("{\"test\":\"test\"}"),
			expectedLoadData: []byte("{\"test\":\"test\"}"),
		},
		{
			name: "Valid GitHub Action with JSON",
			envVars: map[string]string{
				"CI":             "true",
				"GITHUB_ACTIONS": "true",
				"GITHUB_ENV":     "test_env",
				"GITHUB_OUTPUT":  "test_output",
				"RUNNER_TEMP":    "test_temp",
			},
			expectedErr:      nil,
			cmdID:            "123",
			saveData:         mustMarshalYAML(map[string]string{"test": "test"}),
			expectedLoadData: []byte("{\"test\":\"test\"}\n"),
		},

		{
			name: "Valid GitHub Action with JSON",
			envVars: map[string]string{
				"CI":             "true",
				"GITHUB_ACTIONS": "true",
				"GITHUB_ENV":     "test_env",
				"GITHUB_OUTPUT":  "test_output",
				"RUNNER_TEMP":    "test_temp",
			},
			expectedErr:      nil,
			cmdID:            "123",
			saveData:         mustMarshalYAML(buildrc.BuildRC{Version: semver.MustParse("1.0.0"), Packages: []*buildrc.Package{{Name: "test"}}}),
			expectedLoadData: []byte("{\"version\":\"1.0.0\",\"golang\":{\"version\":\"1.20\"},\"packages\":[{\"name\":\"test\"}]}\n"),
		},
		{
			name:        "Not in a GitHub Action",
			envVars:     map[string]string{},
			expectedErr: errors.New("env variable CI is empty"),
		},
		{
			name: "Not in a GitHub Action",
			envVars: map[string]string{
				"CI":             "true",
				"GITHUB_ACTIONS": "false",
			},
			expectedErr: errors.New("not in a github action"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx := context.Background()
			ctx = logging.NewVerboseLogger().WithContext(ctx)

			mockFileAPI := file.NewMemoryFile()

			// Set the environment variables
			for k, v := range tc.envVars {
				os.Setenv(k, v)
				defer func(k string) { os.Unsetenv(k) }(k)
			}

			// Create a new GHActionContentProvider instance
			ghactionCP, err := runner.NewGHActionContentProvider(ctx, mockFileAPI)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			} else {
				assert.NoError(t, err)
			}

			// mockFileAPI.EXPECT().AppendString(ctx, tc.envVars["GITHUB_OUTPUT"], fmt.Sprintf("%s=%s", "result", tc.saveData)).Return(nil)
			// mockFileAPI.EXPECT().Put(ctx, tc.envVars["RUNNER_TEMP"]+"/"+tc.cmdID+".json", tc.saveData).Return(nil)
			// mockFileAPI.EXPECT().Get(ctx, tc.envVars["RUNNER_TEMP"]+"/"+tc.cmdID+".json").Return(tc.expectedLoadData, nil)
			// Mock the file API

			// Implement a simple mock command

			// Save data
			err = ghactionCP.Save(ctx, tc.cmdID, tc.saveData)
			assert.NoError(t, err)

			// Load data
			loadedData, err := ghactionCP.Load(ctx, tc.cmdID)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedLoadData, loadedData)

			// make sure output contains the correct data
			output, err := mockFileAPI.Get(ctx, tc.envVars["GITHUB_OUTPUT"])
			assert.NoError(t, err)

			assert.Contains(t, string(output), "result="+string(tc.saveData))
		})
	}
}
