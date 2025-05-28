package jira

import (
	"fmt"
	goJira "github.com/andygrunwald/go-jira"
	"github.com/tutunak/sigrab/internal/output"
	"github.com/tutunak/sigrab/internal/utils"
)

type Fetcher struct {
	client *goJira.Client
}

func NewFetcher(client *goJira.Client) *Fetcher {
	return &Fetcher{client: client}
}

func (f *Fetcher) FetchBackward(to string, path string) ([]goJira.Issue, error) {
	prefix, endNum, err := utils.ParseIssueKey(to)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issue key %s: %w", to, err)
	}

	var fetchedIssues []goJira.Issue
	for current := endNum; current >= 1; current-- {
		issueKey := fmt.Sprintf("%s-%d", prefix, current)
		issue, err := GetIssue(f.client, issueKey)
		if err != nil {
			// If an issue is not found or there's an error, skip it and continue
			// This matches the test expectation for "skip non-existent issues"
			// Consider logging this error if necessary
			continue
		}

		// Append to maintain fetched order (e.g., N, N-1, N-2) as expected by tests.
		fetchedIssues = append(fetchedIssues, *issue)

		writer := output.NewWriter()
		// writer.WriteToFile expects the directory as the first argument and constructs the filename itself.
		if err := writer.WriteToFile(path, *issue); err != nil {
			// On write error, the test expects nil for the issues slice.
			return nil, fmt.Errorf("failed to write issue %s to directory %s: %w", issueKey, path, err)
		}
	}
	return fetchedIssues, nil
}
