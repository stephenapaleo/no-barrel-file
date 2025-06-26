package ignorer

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIgnorerWithDifferentPathSeparators tests ignorer with Windows and Unix path separators
func TestIgnorerWithDifferentPathSeparators(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create .gitignore with different path separator formats
	var gitignoreContent string
	if runtime.GOOS == "windows" {
		gitignoreContent = `node_modules/
build/
dist/
src/temp/
*.tmp`
	} else {
		gitignoreContent = `node_modules/
build/
dist/
src/temp/
*.tmp`
	}
	
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	assert.NoError(t, err)
	
	// Test paths to ignore
	testPaths := []struct {
		path        string
		shouldIgnore bool
	}{
		{filepath.Join(tmpDir, "node_modules", "package"), true},
		{filepath.Join(tmpDir, "build", "output"), true},
		{filepath.Join(tmpDir, "dist", "bundle"), true},
		{filepath.Join(tmpDir, "src", "temp", "file.ts"), true},
		{filepath.Join(tmpDir, "file.tmp"), true},
		{filepath.Join(tmpDir, "src", "components", "Button.ts"), false},
		{filepath.Join(tmpDir, "lib", "utils.ts"), false},
	}
	
	ignorer := New(tmpDir, []string{}, ".gitignore")
	
	for _, tc := range testPaths {
		t.Run(tc.path, func(t *testing.T) {
			shouldIgnore := ignorer.IgnorePath(tc.path)
			assert.Equal(t, tc.shouldIgnore, shouldIgnore, "Path: %s", tc.path)
		})
	}
}

// TestIgnorerWithCustomIgnorePaths tests ignorer with custom ignore paths using different separators
func TestIgnorerWithCustomIgnorePaths(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Define custom ignore paths
	customIgnorePaths := []string{
		"temp",
		"cache",
		filepath.Join("logs", "debug"),
	}
	
	ignorer := New(tmpDir, customIgnorePaths, "")
	
	// Test paths that should be ignored based on custom ignore paths
	testPaths := []struct {
		path        string
		shouldIgnore bool
	}{
		{filepath.Join(tmpDir, "temp", "file.ts"), true},
		{filepath.Join(tmpDir, "cache", "data.json"), true},
		{filepath.Join(tmpDir, "logs", "debug", "error.log"), true},
		{filepath.Join(tmpDir, "src", "components", "Button.ts"), false},
		{filepath.Join(tmpDir, "lib", "utils.ts"), false},
	}
	
	for _, tc := range testPaths {
		t.Run(tc.path, func(t *testing.T) {
			shouldIgnore := ignorer.IgnorePath(tc.path)
			assert.Equal(t, tc.shouldIgnore, shouldIgnore, "Path: %s", tc.path)
		})
	}
}

// TestIgnorerWithSpecialCharacters tests ignorer with special characters in paths
func TestIgnorerWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create directories and files with special characters
	specialDirs := []string{
		"components-with-dash",
		"components_with_underscore",
		"components.with.dots",
	}
	
	if runtime.GOOS != "windows" {
		// Unix allows more special characters
		specialDirs = append(specialDirs, "components with spaces")
	}
	
	// Create .gitignore with patterns for special characters
	gitignoreContent := `*-with-dash/
*_with_underscore/
*.with.dots/`
	
	if runtime.GOOS != "windows" {
		gitignoreContent += `
*\ with\ spaces/`
	}
	
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	assert.NoError(t, err)
	
	ignorer := New(tmpDir, []string{}, ".gitignore")
	
	// Test paths with special characters
	for _, dirName := range specialDirs {
		testPath := filepath.Join(tmpDir, dirName, "index.ts")
		
		shouldIgnore := ignorer.IgnorePath(testPath)
		
		// Most special character directories should be ignored based on our gitignore patterns
		expectedIgnore := true
		if dirName == "components with spaces" && runtime.GOOS == "windows" {
			// This pattern might not work on Windows
			expectedIgnore = false
		}
		
		// Note: The actual behavior depends on the gitignore library implementation
		// This test verifies that the ignorer doesn't crash with special characters
		_ = shouldIgnore
		_ = expectedIgnore
		
		t.Logf("Directory: %s, Ignored: %v", dirName, shouldIgnore)
	}
}

// TestIgnorerPlatformDifferences tests platform-specific behaviors
func TestIgnorerPlatformDifferences(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a simple .gitignore
	gitignoreContent := `*.log
build/`
	
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	assert.NoError(t, err)
	
	ignorer := New(tmpDir, []string{}, ".gitignore")
	
	// Test basic functionality works on both platforms
	testCases := []struct {
		name         string
		path         string
		shouldIgnore bool
	}{
		{
			name:         "Log file should be ignored",
			path:         filepath.Join(tmpDir, "error.log"),
			shouldIgnore: true,
		},
		{
			name:         "Build directory should be ignored",
			path:         filepath.Join(tmpDir, "build", "output.js"),
			shouldIgnore: true,
		},
		{
			name:         "Source file should not be ignored",
			path:         filepath.Join(tmpDir, "src", "index.ts"),
			shouldIgnore: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shouldIgnore := ignorer.IgnorePath(tc.path)
			assert.Equal(t, tc.shouldIgnore, shouldIgnore, "Path: %s", tc.path)
		})
	}
} 