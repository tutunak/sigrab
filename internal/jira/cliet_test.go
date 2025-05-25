package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	goJira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
)

func TestGetIssue_HTTPMock_Success(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/rest/api/2/issue/PROJ-123")

		// Return mock response
		issue := goJira.Issue{
			Key: "PROJ-123",
			Fields: &goJira.IssueFields{
				Summary: "Test Issue",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issue)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClient("testuser", "testtoken", server.URL)

	// Test the function
	result, err := GetIssue(client, "PROJ-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PROJ-123", result.Key)
	assert.Equal(t, "Test Issue", result.Fields.Summary)
}

func TestGetIssue_HTTPMock_Error(t *testing.T) {
	// Create a mock HTTP server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errorMessages":["Issue does not exist"]}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClient("testuser", "testtoken", server.URL)

	// Test the function
	result, err := GetIssue(client, "PROJ-404")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}
