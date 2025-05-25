package jira

import (
	goJira "github.com/andygrunwald/go-jira"
)

func NewClient(username string, token string, url string) *goJira.Client {
	tp := goJira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	jiraClient, err := goJira.NewClient(tp.Client(), url)
	if err != nil {
		panic(err)
	}

	return jiraClient
}

func GetIssue(client *goJira.Client, issueKey string) (*goJira.Issue, error) {
	issue, _, err := client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}
