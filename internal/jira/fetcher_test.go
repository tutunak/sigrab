package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	goJira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetcher_FetchBackward(t *testing.T) {
	t.Run("successful fetch multiple issues", func(t *testing.T) {
		// Create temp directory for output
		tempDir, err := os.MkdirTemp("", "fetcher_test_*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Mock server to return 3 issues
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			issueKey := filepath.Base(r.URL.Path)
			switch issueKey {
			case "PROJ-1", "PROJ-2", "PROJ-3":
				issue := goJira.Issue{
					Key: issueKey,
					Fields: &goJira.IssueFields{
						Summary: fmt.Sprintf("Test Issue %s", issueKey),
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(issue)
			default:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"errorMessages":["Issue does not exist"]}`))
			}
		}))
		defer server.Close()

		// Create client and fetcher
		client := NewClient("user", "token", server.URL)
		fetcher := NewFetcher(client)

		// Test the function
		issues, err := fetcher.FetchBackward("PROJ-3", tempDir)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, issues, 3)
		assert.Equal(t, "PROJ-1", issues[2].Key)
		assert.Equal(t, "PROJ-2", issues[1].Key)
		assert.Equal(t, "PROJ-3", issues[0].Key)

		// Verify files were created
		for i := 1; i <= 3; i++ {
			filePath := filepath.Join(tempDir, fmt.Sprintf("PROJ-%d.json", i))
			_, err := os.Stat(filePath)
			assert.NoError(t, err, "File should exist: %s", filePath)
		}
	})

	t.Run("skip non-existent issues", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "fetcher_test_*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Mock server to return only issues 1 and 3
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			issueKey := filepath.Base(r.URL.Path)
			switch issueKey {
			case "PROJ-1", "PROJ-3":
				issue := goJira.Issue{
					Key: issueKey,
					Fields: &goJira.IssueFields{
						Summary: fmt.Sprintf("Test Issue %s", issueKey),
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(issue)
			default:
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"errorMessages":["Issue does not exist"]}`))
			}
		}))
		defer server.Close()

		client := NewClient("user", "token", server.URL)
		fetcher := NewFetcher(client)

		issues, err := fetcher.FetchBackward("PROJ-3", tempDir)

		assert.NoError(t, err)
		assert.Len(t, issues, 2)
		assert.Equal(t, "PROJ-1", issues[1].Key)
		assert.Equal(t, "PROJ-3", issues[0].Key)

		// Verify only files 1 and 3 exist
		file1 := filepath.Join(tempDir, "PROJ-1.json")
		file2 := filepath.Join(tempDir, "PROJ-2.json")
		file3 := filepath.Join(tempDir, "PROJ-3.json")

		_, err = os.Stat(file1)
		assert.NoError(t, err)
		_, err = os.Stat(file2)
		assert.True(t, os.IsNotExist(err))
		_, err = os.Stat(file3)
		assert.NoError(t, err)
	})

	t.Run("invalid issue key", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "fetcher_test_*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		client := NewClient("user", "token", "http://example.com")
		fetcher := NewFetcher(client)

		// Test with invalid key format
		issues, err := fetcher.FetchBackward("INVALID_KEY", tempDir)
		assert.Error(t, err)
		assert.Nil(t, issues)
	})

	t.Run("write error", func(t *testing.T) {
		// Use a non-existent directory to cause write error
		nonExistentDir := filepath.Join(os.TempDir(), "non_existent_dir_"+t.Name())

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			issue := goJira.Issue{
				Key: "PROJ-1",
				Fields: &goJira.IssueFields{
					Summary: "Test Issue",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(issue)
		}))
		defer server.Close()

		client := NewClient("user", "token", server.URL)
		fetcher := NewFetcher(client)

		issues, err := fetcher.FetchBackward("PROJ-1", nonExistentDir)
		assert.Error(t, err)
		assert.Nil(t, issues)
	})

	t.Run("no issues found", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "fetcher_test_*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Mock server to return 404 for all issues
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"errorMessages":["Issue does not exist"]}`))
		}))
		defer server.Close()

		client := NewClient("user", "token", server.URL)
		fetcher := NewFetcher(client)

		issues, err := fetcher.FetchBackward("PROJ-5", tempDir)
		assert.NoError(t, err)
		assert.Empty(t, issues)
	})
}
