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
