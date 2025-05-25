package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedError  string
		expectedConfig *Config
	}{
		{
			name: "successful config load with all env vars set",
			envVars: map[string]string{
				"JIRA_API_TOKEN":  "test-token-123",
				"JIRA_USER_EMAIL": "test@example.com",
			},
			expectedError: "",
			expectedConfig: &Config{
				APIToken:  "test-token-123",
				UserEmail: "test@example.com",
			},
		},
		{
			name: "missing JIRA_API_TOKEN",
			envVars: map[string]string{
				"JIRA_USER_EMAIL": "test@example.com",
			},
			expectedError:  "missing required environment variables: JIRA_API_TOKEN",
			expectedConfig: nil,
		},
		{
			name: "missing JIRA_USER_EMAIL",
			envVars: map[string]string{
				"JIRA_API_TOKEN": "test-token-123",
			},
			expectedError:  "missing required environment variables: JIRA_USER_EMAIL",
			expectedConfig: nil,
		},
		{
			name:           "missing both environment variables",
			envVars:        map[string]string{},
			expectedError:  "missing required environment variables: JIRA_API_TOKEN, JIRA_USER_EMAIL",
			expectedConfig: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: Clear and set environment variables for this test
			clearTestEnvVars()
			setTestEnvVars(tt.envVars)

			// Cleanup after test
			defer clearTestEnvVars()

			// Execute
			config, err := LoadConfig()

			// Assert
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, config)
			} else {
				require.NoError(t, err)
				require.NotNil(t, config)
				assert.Equal(t, tt.expectedConfig.APIToken, config.APIToken)
				assert.Equal(t, tt.expectedConfig.UserEmail, config.UserEmail)
			}
		})
	}
}

func TestLoadConfig_ErrorMessageFormat(t *testing.T) {
	clearTestEnvVars()
	defer clearTestEnvVars()

	config, err := LoadConfig()

	require.Error(t, err)
	assert.Nil(t, config)

	errMsg := err.Error()
	assert.Contains(t, errMsg, "missing required environment variables:")
	assert.Contains(t, errMsg, "JIRA_API_TOKEN")
	assert.Contains(t, errMsg, "JIRA_USER_EMAIL")
	assert.Contains(t, errMsg, "JIRA_API_TOKEN, JIRA_USER_EMAIL")
}

func BenchmarkLoadConfig(b *testing.B) {
	err := os.Setenv("JIRA_API_TOKEN", "benchmark-token")
	require.NoError(b, err)
	err = os.Setenv("JIRA_USER_EMAIL", "benchmark@example.com")
	require.NoError(b, err)

	defer func() {
		os.Unsetenv("JIRA_API_TOKEN")
		os.Unsetenv("JIRA_USER_EMAIL")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadConfig()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Helper functions for test setup and cleanup
func clearTestEnvVars() {
	os.Unsetenv("JIRA_API_TOKEN")
	os.Unsetenv("JIRA_USER_EMAIL")
}

func setTestEnvVars(vars map[string]string) {
	for key, value := range vars {
		os.Setenv(key, value)
	}
}
