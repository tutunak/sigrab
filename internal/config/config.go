package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APIToken  string
	UserEmail string
}

func LoadConfig() (*Config, error) {
	var missingVars []string

	jiraToken := os.Getenv("JIRA_API_TOKEN")
	if jiraToken == "" {
		missingVars = append(missingVars, "JIRA_API_TOKEN")
	}

	userEmail := os.Getenv("JIRA_USER_EMAIL")
	if userEmail == "" {
		missingVars = append(missingVars, "JIRA_USER_EMAIL")
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s",
			strings.Join(missingVars, ", "))
	}

	return &Config{
		APIToken:  jiraToken,
		UserEmail: userEmail,
	}, nil
}
