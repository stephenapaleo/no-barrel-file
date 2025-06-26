package tests

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPathNormalization tests path normalization across platforms
func TestPathNormalization(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedUnix   string
		expectedWindows string
	}{
		{
			name:           "Unix path with forward slashes",
			input:          "src/components/Button/index.ts",
			expectedUnix:   "src/components/Button/index.ts",
			expectedWindows: "src\\components\\Button\\index.ts",
		},
		{
			name:           "Mixed path separators",
			input:          "src\\components/Button\\index.ts",
			expectedUnix:   "src/components/Button/index.ts",
			expectedWindows: "src\\components\\Button\\index.ts",
		},
		{
			name:           "Absolute Unix path",
			input:          "/home/user/project/src/index.ts",
			expectedUnix:   "/home/user/project/src/index.ts",
			expectedWindows: "\\home\\user\\project\\src\\index.ts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := filepath.Clean(tt.input)
			
			if runtime.GOOS == "windows" {
				// On Windows, filepath.Clean normalizes to Windows separators
				assert.Equal(t, tt.expectedWindows, normalized)
			} else {
				// On Unix systems, filepath.Clean normalizes to Unix separators
				assert.Equal(t, tt.expectedUnix, normalized)
			}
		})
	}
}

// TestPathJoin tests cross-platform path joining
func TestPathJoin(t *testing.T) {
	tests := []struct {
		name           string
		parts          []string
		expectedUnix   string
		expectedWindows string
	}{
		{
			name:           "Basic path join",
			parts:          []string{"src", "components", "Button", "index.ts"},
			expectedUnix:   "src/components/Button/index.ts",
			expectedWindows: "src\\components\\Button\\index.ts",
		},
		{
			name:           "Path join with empty parts",
			parts:          []string{"src", "", "components", "index.ts"},
			expectedUnix:   "src/components/index.ts",
			expectedWindows: "src\\components\\index.ts",
		},
		{
			name:           "Path join with current directory",
			parts:          []string{".", "src", "index.ts"},
			expectedUnix:   "src/index.ts",
			expectedWindows: "src\\index.ts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filepath.Join(tt.parts...)
			
			if runtime.GOOS == "windows" {
				assert.Equal(t, tt.expectedWindows, result)
			} else {
				assert.Equal(t, tt.expectedUnix, result)
			}
		})
	}
}

// TestPathMatching tests case sensitivity differences
func TestPathMatching(t *testing.T) {
	tests := []struct {
		name          string
		path1         string
		path2         string
		shouldMatchWindows bool
		shouldMatchUnix    bool
	}{
		{
			name:          "Exact match",
			path1:         "src/Index.ts",
			path2:         "src/Index.ts",
			shouldMatchWindows: true,
			shouldMatchUnix:    true,
		},
		{
			name:          "Case difference",
			path1:         "src/Index.ts",
			path2:         "src/index.ts",
			shouldMatchWindows: true,  // Windows is case-insensitive
			shouldMatchUnix:    false, // Unix is case-sensitive
		},
		{
			name:          "Directory case difference",
			path1:         "SRC/components/index.ts",
			path2:         "src/components/index.ts",
			shouldMatchWindows: true,  // Windows is case-insensitive
			shouldMatchUnix:    false, // Unix is case-sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normalize paths for comparison
			normalized1 := filepath.Clean(tt.path1)
			normalized2 := filepath.Clean(tt.path2)
			
			if runtime.GOOS == "windows" {
				// On Windows, compare case-insensitively
				matches := strings.EqualFold(normalized1, normalized2)
				assert.Equal(t, tt.shouldMatchWindows, matches)
			} else {
				// On Unix, compare case-sensitively
				matches := normalized1 == normalized2
				assert.Equal(t, tt.shouldMatchUnix, matches)
			}
		})
	}
}

// TestAbsolutePaths tests absolute path handling
func TestAbsolutePaths(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		isAbsolute bool
	}{
		{
			name:     "Unix absolute path",
			path:     "/home/user/project",
			isAbsolute: runtime.GOOS != "windows", // Only absolute on Unix
		},
		{
			name:     "Windows absolute path with drive",
			path:     "C:\\Users\\user\\project",
			isAbsolute: runtime.GOOS == "windows", // Only absolute on Windows
		},
		{
			name:     "Relative path",
			path:     "src/components",
			isAbsolute: false,
		},
		{
			name:     "Current directory",
			path:     ".",
			isAbsolute: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isAbs := filepath.IsAbs(tt.path)
			assert.Equal(t, tt.isAbsolute, isAbs)
		})
	}
}

// TestRelativePaths tests relative path calculation
func TestRelativePaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Run("Windows relative paths", func(t *testing.T) {
			base := "C:\\Users\\user\\project"
			target := "C:\\Users\\user\\project\\src\\components\\Button\\index.ts"
			
			rel, err := filepath.Rel(base, target)
			assert.NoError(t, err)
			assert.Equal(t, "src\\components\\Button\\index.ts", rel)
		})
		
		t.Run("Windows cross-drive relative paths", func(t *testing.T) {
			base := "C:\\Users\\user\\project"
			target := "D:\\other\\project\\src\\index.ts"
			
			rel, err := filepath.Rel(base, target)
			// Cross-drive paths should return an error on Windows
			assert.Error(t, err, "Cross-drive paths should return an error")
			assert.Empty(t, rel, "Relative path should be empty when error occurs")
		})
	} else {
		t.Run("Unix relative paths", func(t *testing.T) {
			base := "/home/user/project"
			target := "/home/user/project/src/components/Button/index.ts"
			
			rel, err := filepath.Rel(base, target)
			assert.NoError(t, err)
			assert.Equal(t, "src/components/Button/index.ts", rel)
		})
		
		t.Run("Unix relative paths with parent directory", func(t *testing.T) {
			base := "/home/user/project/src"
			target := "/home/user/project/lib/utils.ts"
			
			rel, err := filepath.Rel(base, target)
			assert.NoError(t, err)
			assert.Equal(t, "../lib/utils.ts", rel)
		})
	}
}

// TestPathSeparatorConversion tests path separator conversion
func TestPathSeparatorConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unix path to slash",
			input:    "src/components/Button/index.ts",
			expected: "src/components/Button/index.ts",
		},
		{
			name:     "Windows path to slash",
			input:    "src\\components\\Button\\index.ts",
			expected: "src/components/Button/index.ts",
		},
		{
			name:     "Mixed separators to slash",
			input:    "src/components\\Button/index.ts",
			expected: "src/components/Button/index.ts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filepath.ToSlash(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFileSystemOperations tests file system operations across platforms
func TestFileSystemOperations(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Test creating directories with different path separators
	testPaths := []string{
		filepath.Join("src", "components"),
		filepath.Join("lib", "utils"),
		filepath.Join("test", "data", "input"),
	}
	
	for _, testPath := range testPaths {
		fullPath := filepath.Join(tmpDir, testPath)
		err := os.MkdirAll(fullPath, 0755)
		assert.NoError(t, err)
		
		// Verify directory exists
		stat, err := os.Stat(fullPath)
		assert.NoError(t, err)
		assert.True(t, stat.IsDir())
	}
	
	// Test creating files with different path separators
	testFiles := []string{
		filepath.Join("src", "components", "index.ts"),
		filepath.Join("lib", "utils", "helper.ts"),
		filepath.Join("test", "data", "input", "barrel.ts"),
	}
	
	for _, testFile := range testFiles {
		fullPath := filepath.Join(tmpDir, testFile)
		// Create parent directory if it doesn't exist
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		assert.NoError(t, err)
		
		// Create file
		file, err := os.Create(fullPath)
		assert.NoError(t, err)
		file.Close()
		
		// Verify file exists
		stat, err := os.Stat(fullPath)
		assert.NoError(t, err)
		assert.False(t, stat.IsDir())
	}
}

// TestPathPattern tests path pattern matching (for gitignore-like functionality)
func TestPathPattern(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		path        string
		shouldMatch bool
	}{
		{
			name:        "Exact match",
			pattern:     "node_modules",
			path:        "node_modules",
			shouldMatch: true,
		},
		{
			name:        "Pattern with slash",
			pattern:     "src/components",
			path:        filepath.Join("src", "components"),
			shouldMatch: true,
		},
		{
			name:        "Pattern with wildcard",
			pattern:     "*.ts",
			path:        "index.ts",
			shouldMatch: true,
		},
		{
			name:        "Pattern with wildcard no match",
			pattern:     "*.ts",
			path:        "index.js",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Normalize both pattern and path for comparison
			normalizedPattern := filepath.ToSlash(tt.pattern)
			normalizedPath := filepath.ToSlash(tt.path)
			
			// Simple pattern matching (would be more complex in real implementation)
			var matches bool
			if strings.Contains(normalizedPattern, "*") {
				// Simple wildcard matching
				if strings.HasSuffix(normalizedPattern, "*") {
					prefix := strings.TrimSuffix(normalizedPattern, "*")
					matches = strings.HasPrefix(normalizedPath, prefix)
				} else if strings.HasPrefix(normalizedPattern, "*") {
					suffix := strings.TrimPrefix(normalizedPattern, "*")
					matches = strings.HasSuffix(normalizedPath, suffix)
				}
			} else {
				matches = normalizedPattern == normalizedPath
			}
			
			assert.Equal(t, tt.shouldMatch, matches)
		})
	}
} 