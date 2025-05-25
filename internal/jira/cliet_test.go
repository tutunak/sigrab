package jira

import (
	"errors"
	"testing"

	goJira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockJiraClient is a mock implementation for testing
type MockJiraClient struct {
	mock.Mock
	Issue *MockIssueService
}

type MockIssueService struct {
	mock.Mock
}

func (m *MockIssueService) Get(issueKey string, options *goJira.GetQueryOptions) (*goJira.Issue, *goJira.Response, error) {
	args := m.Called(issueKey, options)
	return args.Get(0).(*goJira.Issue), args.Get(1).(*goJira.Response), args.Error(2)
}

// JiraTestSuite defines the test suite
type JiraTestSuite struct {
	suite.Suite
	mockClient *MockJiraClient
}

func (suite *JiraTestSuite) SetupTest() {
	suite.mockClient = &MockJiraClient{
		Issue: &MockIssueService{},
	}
}

func TestJiraTestSuite(t *testing.T) {
	suite.Run(t, new(JiraTestSuite))
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		username string
		token    string
		url      string
		wantErr  bool
	}{
		{
			name:     "Valid credentials and URL",
			username: "testuser",
			token:    "testtoken",
			url:      "https://company.atlassian.net",
			wantErr:  false,
		},
		{
			name:     "Empty username",
			username: "",
			token:    "testtoken",
			url:      "https://company.atlassian.net",
			wantErr:  false, // BasicAuth allows empty username
		},
		{
			name:     "Empty token",
			username: "testuser",
			token:    "",
			url:      "https://company.atlassian.net",
			wantErr:  false, // BasicAuth allows empty password
		},
		{
			name:     "Empty URL",
			username: "testuser",
			token:    "testtoken",
			url:      "",
			wantErr:  false, // Will use default base URL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Panics(t, func() {
					NewClient(tt.username, tt.token, tt.url)
				})
			} else {
				assert.NotPanics(t, func() {
					client := NewClient(tt.username, tt.token, tt.url)
					assert.NotNil(t, client)
				})
			}
		})
	}
}

func TestNewClient_InvalidURL(t *testing.T) {
	// Test with an invalid URL that would cause goJira.NewClient to fail
	assert.Panics(t, func() {
		NewClient("user", "token", "://invalid-url")
	})
}

func (suite *JiraTestSuite) TestGetIssue_Success() {
	// Arrange
	issueKey := "PROJ-123"
	expectedIssue := &goJira.Issue{
		Key: issueKey,
		Fields: &goJira.IssueFields{
			Summary: "Test Issue",
		},
	}
	response := &goJira.Response{}

	// Create a real client for this test
	client := &goJira.Client{
		Issue: suite.mockClient.Issue,
	}

	suite.mockClient.Issue.On("Get", issueKey, (*goJira.GetQueryOptions)(nil)).Return(expectedIssue, response, nil)

	// Act
	result, err := GetIssue(client, issueKey)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), issueKey, result.Key)
	assert.Equal(suite.T(), "Test Issue", result.Fields.Summary)
	suite.mockClient.Issue.AssertExpectations(suite.T())
}

func (suite *JiraTestSuite) TestGetIssue_Error() {
	// Arrange
	issueKey := "PROJ-404"
	expectedError := errors.New("issue not found")
	response := &goJira.Response{}

	client := &goJira.Client{
		Issue: suite.mockClient.Issue,
	}

	suite.mockClient.Issue.On("Get", issueKey, (*goJira.GetQueryOptions)(nil)).Return((*goJira.Issue)(nil), response, expectedError)

	// Act
	result, err := GetIssue(client, issueKey)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockClient.Issue.AssertExpectations(suite.T())
}

func TestGetIssue_EmptyIssueKey(t *testing.T) {
	// Test with empty issue key
	client := &goJira.Client{}

	result, err := GetIssue(client, "")

	// The actual behavior depends on the go-jira library implementation
	// This test documents the current behavior
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetIssue_NilClient(t *testing.T) {
	// Test with nil client - should panic
	assert.Panics(t, func() {
		GetIssue(nil, "PROJ-123")
	})
}

func TestIntegration_NewClientAndGetIssue(t *testing.T) {
	// Skip this test in unit test runs
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// This would be an integration test with a real JIRA instance
	// You'd need valid credentials and a test JIRA instance
	t.Skip("Integration test - requires real JIRA instance")

	username := "test@example.com"
	token := "your-api-token"
	url := "https://your-domain.atlassian.net"

	client := NewClient(username, token, url)
	assert.NotNil(t, client)

	// Test with a known issue key in your test instance
	issue, err := GetIssue(client, "TEST-1")

	if err != nil {
		t.Logf("Expected error for test issue: %v", err)
	} else {
		assert.NotNil(t, issue)
		assert.NotEmpty(t, issue.Key)
	}
}
