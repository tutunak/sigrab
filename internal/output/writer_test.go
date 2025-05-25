package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	goJira "github.com/andygrunwald/go-jira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriter_WriteToFile(t *testing.T) {
	// Create a test issue
	testIssue := goJira.Issue{
		Key: "TEST-123",
		Fields: &goJira.IssueFields{
			Summary:     "Test Issue",
			Description: "This is a test issue",
		},
	}

	t.Run("success case", func(t *testing.T) {
		// Setup temp directory
		tempDir, err := os.MkdirTemp("", "writer_test_*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create writer and write issue
		writer := NewWriter()
		err = writer.WriteToFile(tempDir, testIssue)
		assert.NoError(t, err)

		// Verify file exists
		expectedPath := filepath.Join(tempDir, "TEST-123.json")
		_, err = os.Stat(expectedPath)
		assert.NoError(t, err)

		// Verify file content
		fileContent, err := os.ReadFile(expectedPath)
		require.NoError(t, err)

		var parsedIssue goJira.Issue
		err = json.Unmarshal(fileContent, &parsedIssue)
		require.NoError(t, err)

		assert.Equal(t, testIssue.Key, parsedIssue.Key)
		assert.Equal(t, testIssue.Fields.Summary, parsedIssue.Fields.Summary)
		assert.Equal(t, testIssue.Fields.Description, parsedIssue.Fields.Description)
	})

	t.Run("path doesn't exist", func(t *testing.T) {
		// Setup a path that doesn't exist
		nonExistentPath := filepath.Join(os.TempDir(), "non_existent_dir_"+t.Name())

		// Ensure the path doesn't exist
		_, err := os.Stat(nonExistentPath)
		require.True(t, os.IsNotExist(err))

		// Create writer and attempt to write issue
		writer := NewWriter()
		err = writer.WriteToFile(nonExistentPath, testIssue)
		assert.Error(t, err)
	})

	t.Run("path is a file, not a directory", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "writer_test_file_*")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Create writer and attempt to write issue
		writer := NewWriter()
		err = writer.WriteToFile(tempFile.Name(), testIssue)
		assert.Error(t, err)
	})

	t.Run("no write permission", func(t *testing.T) {
		// Skip if running as root, as root can write anywhere
		if os.Geteuid() == 0 {
			t.Skip("Test skipped when running as root")
		}

		// Create temp directory with no write permissions
		tempDir, err := os.MkdirTemp("", "writer_test_noperm_*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Change permissions to read-only
		err = os.Chmod(tempDir, 0500) // r-x------
		require.NoError(t, err)

		// Create writer and attempt to write issue
		writer := NewWriter()
		err = writer.WriteToFile(tempDir, testIssue)
		assert.Error(t, err)
	})

	t.Run("special characters in issue key", func(t *testing.T) {
		// Setup temp directory
		tempDir, err := os.MkdirTemp("", "writer_test_*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create an issue with special characters in key
		specialIssue := testIssue
		specialIssue.Key = "SPECIAL-456"

		// Create writer and write issue
		writer := NewWriter()
		err = writer.WriteToFile(tempDir, specialIssue)
		assert.NoError(t, err)

		// Verify file exists
		expectedPath := filepath.Join(tempDir, "SPECIAL-456.json")
		_, err = os.Stat(expectedPath)
		assert.NoError(t, err)
	})
}
