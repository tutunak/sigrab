package output

import (
	"encoding/json"
	"fmt"
	goJira "github.com/andygrunwald/go-jira"
	"os"
)

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) WriteToFile(path string, issue goJira.Issue) error {
	file := fmt.Sprintf("%s/%s.json", path, issue.Key)
	data, err := json.MarshalIndent(issue, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal issue to JSON: %w", err)
	}

	// Write JSON data to file
	if err := os.WriteFile(file, data, 0644); err != nil {
		return fmt.Errorf("failed to write issue to file %s: %w", file, err)
	}

	return nil
}
