package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseIssueKey(issueKey string) (prefix string, number int, err error) {
	parts := strings.Split(issueKey, "-")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid issue key format: %s", issueKey)
	}

	if parts[0] == "" || parts[1] == "" {
		return "", 0, fmt.Errorf("issue key prefix or number cannot be empty: %s", issueKey)
	}

	prefix = parts[0]
	number, err = strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid issue number: %s", parts[1])
	}

	return prefix, number, nil
}
