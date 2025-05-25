package utils

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampInitialization(t *testing.T) {
	// Since Timestamp is initialized at package import time,
	// we can only verify it has the correct format
	_, err := time.Parse("20060102_150405", Timestamp)
	assert.NoError(t, err, "Timestamp should be in the format 20060102_150405")
}

func TestInitDir(t *testing.T) {
	// Setup temporary directory for testing
	testDir, err := os.MkdirTemp("", "sigrab_test_*")
	assert.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Failed to clean up test directory: %v", err)
		}
	}(testDir) // Cleanup after test

	tests := []struct {
		name        string
		inputPath   string
		expectPanic bool
		check       func(t *testing.T, resultPath string)
	}{
		{
			name:        "empty path should use default",
			inputPath:   "",
			expectPanic: false,
			check: func(t *testing.T, resultPath string) {
				assert.Equal(t, filepath.Join("/tmp/sigrab", Timestamp), resultPath)
				_, err := os.Stat(resultPath)
				assert.NoError(t, err, "Directory should exist")
			},
		},
		{
			name:        "valid path should create directory",
			inputPath:   testDir,
			expectPanic: false,
			check: func(t *testing.T, resultPath string) {
				assert.Equal(t, filepath.Join(testDir, Timestamp), resultPath)
				_, err := os.Stat(resultPath)
				assert.NoError(t, err, "Directory should exist")
			},
		},
		{
			name:        "non-writable path should panic",
			inputPath:   "/root/nonexistent/directory", // Most users won't have write permission here
			expectPanic: true,
			check:       func(t *testing.T, resultPath string) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					InitDir(tt.inputPath)
				})
			} else {
				var resultPath string
				assert.NotPanics(t, func() {
					resultPath = InitDir(tt.inputPath)
				})
				tt.check(t, resultPath)
			}
		})
	}
}

func TestInitDir_CleanupBetweenTests(t *testing.T) {
	// Setup
	testDir, err := os.MkdirTemp("", "sigrab_test_*")
	assert.NoError(t, err)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Failed to clean up test directory: %v", err)
		}
	}(testDir)

	// Run multiple times to ensure no collision
	resultPath1 := InitDir(testDir)
	time.Sleep(1 * time.Second) // Ensure timestamp changes

	// Save current timestamp
	oldTimestamp := Timestamp
	// Update timestamp for second test
	Timestamp = time.Now().Format("20060102_150405")

	resultPath2 := InitDir(testDir)

	// Verify paths are different
	assert.NotEqual(t, resultPath1, resultPath2, "Directory paths should be different due to timestamp")

	// Verify both directories exist
	_, err = os.Stat(resultPath1)
	assert.NoError(t, err)
	_, err = os.Stat(resultPath2)
	assert.NoError(t, err)

	// Restore original timestamp for other tests
	Timestamp = oldTimestamp
}
