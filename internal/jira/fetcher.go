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

func (f *Fetcher) FetchBackward(to string, path string) error {
	prefix, endNum, err := utils.ParseIssueKey(to)
	if err != nil {

		return fmt.Errorf("failed to parse issue key %s: %w", to, err)
	}

	for current := endNum; current >= 1; current-- {
		issueKey := fmt.Sprintf("%s-%d", prefix, current)
		issue, err := GetIssue(f.client, issueKey)
		if err != nil {
			continue
		}

		writer := output.NewWriter()
		if err := writer.WriteToFile(path, *issue); err != nil {

			return fmt.Errorf("failed to write issue %s to file: %w", issueKey, err)
		}

		// Prepend to maintain order
	}
	return nil
}
