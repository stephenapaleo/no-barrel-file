package resolver

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResolverWithWindowsPaths tests resolver with Windows-style paths
func TestResolverWithWindowsPaths(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create tsconfig.json with Windows-style paths in content
	var tsconfigContent string
	if runtime.GOOS == "windows" {
		// On Windows, we can test with both forward and back slashes
		tsconfigContent = `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@components/*": ["src/components/*"],
      "@utils/*": ["src\\utils\\*"]
    }
  }
}`
	} else {
		// On Unix, we normalize to forward slashes
		tsconfigContent = `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@components/*": ["src/components/*"],
      "@utils/*": ["src/utils/*"]
    }
  }
}`
	}
	
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err := os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Create directory structure
	srcDir := filepath.Join(tmpDir, "src")
	componentsDir := filepath.Join(srcDir, "components")
	utilsDir := filepath.Join(srcDir, "utils")
	
	err = os.MkdirAll(componentsDir, 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(utilsDir, 0755)
	assert.NoError(t, err)
	
	// Test resolver
	tsconfigRelPath := "tsconfig.json"
	resolver := New(tmpDir, &tsconfigRelPath)
	
	// Test alias resolution for components path
	componentsAlias := resolver.AliasPath(componentsDir)
	assert.Equal(t, "@components", componentsAlias.ShortPath)
	assert.Contains(t, componentsAlias.FullPath, "@components")
	
	// Test alias resolution for utils path
	utilsAlias := resolver.AliasPath(utilsDir)
	assert.Equal(t, "@utils", utilsAlias.ShortPath)
	assert.Contains(t, utilsAlias.FullPath, "@utils")
}

// TestResolverWithRelativePaths tests resolver with different relative path formats
func TestResolverWithRelativePaths(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create nested directory for tsconfig
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)
	
	// Create tsconfig.json with relative baseUrl
	tsconfigContent := `{
  "compilerOptions": {
    "baseUrl": "../src",
    "paths": {
      "@/*": ["./*"],
      "@components/*": ["components/*"]
    }
  }
}`
	
	tsconfigPath := filepath.Join(configDir, "tsconfig.json")
	err = os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Create src directory structure
	srcDir := filepath.Join(tmpDir, "src")
	componentsDir := filepath.Join(srcDir, "components")
	err = os.MkdirAll(componentsDir, 0755)
	assert.NoError(t, err)
	
	// Test resolver with relative tsconfig path
	tsconfigRelPath := "config/tsconfig.json"
	resolver := New(tmpDir, &tsconfigRelPath)
	
	// Test alias resolution
	componentsAlias := resolver.AliasPath(componentsDir)
	assert.Equal(t, "@components", componentsAlias.ShortPath)
}

// TestResolverCaseSensitivity tests resolver with case-sensitive vs case-insensitive paths
func TestResolverCaseSensitivity(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create tsconfig.json
	tsconfigContent := `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@Components/*": ["src/Components/*"],
      "@components/*": ["src/components/*"]
    }
  }
}`
	
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err := os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Create directory structure with different cases
	srcDir := filepath.Join(tmpDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	assert.NoError(t, err)
	
	// Try to create both cases (may fail on case-insensitive file systems)
	componentsLowerDir := filepath.Join(srcDir, "components")
	componentsUpperDir := filepath.Join(srcDir, "Components")
	
	err = os.MkdirAll(componentsLowerDir, 0755)
	assert.NoError(t, err)
	
	err = os.MkdirAll(componentsUpperDir, 0755)
	if err != nil && runtime.GOOS == "windows" {
		// On Windows, this might fail due to case insensitivity
		t.Logf("Could not create both case variants on case-insensitive filesystem: %v", err)
	}
	
	// Test resolver
	tsconfigRelPath := "tsconfig.json"
	resolver := New(tmpDir, &tsconfigRelPath)
	
	// Test alias resolution for lowercase
	lowerAlias := resolver.AliasPath(componentsLowerDir)
	
	if runtime.GOOS == "windows" {
		// On Windows, should match case-insensitively
		assert.True(t, lowerAlias.ShortPath == "@components" || lowerAlias.ShortPath == "@Components")
	} else {
		// On Unix, should match case-sensitively
		assert.Equal(t, "@components", lowerAlias.ShortPath)
	}
}

// TestResolverWithAbsolutePaths tests resolver with absolute paths in tsconfig
func TestResolverWithAbsolutePaths(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create absolute path for baseUrl
	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	assert.NoError(t, err)
	
	// Create tsconfig.json with absolute baseUrl
	var tsconfigContent string
	if runtime.GOOS == "windows" {
		// Windows absolute path
		tsconfigContent = `{
  "compilerOptions": {
    "baseUrl": "` + strings.ReplaceAll(srcDir, "\\", "\\\\") + `",
    "paths": {
      "@/*": ["./*"]
    }
  }
}`
	} else {
		// Unix absolute path
		tsconfigContent = `{
  "compilerOptions": {
    "baseUrl": "` + srcDir + `",
    "paths": {
      "@/*": ["./*"]
    }
  }
}`
	}
	
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err = os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Create components directory
	componentsDir := filepath.Join(srcDir, "components")
	err = os.MkdirAll(componentsDir, 0755)
	assert.NoError(t, err)
	
	// Test resolver
	tsconfigRelPath := "tsconfig.json"
	resolver := New(tmpDir, &tsconfigRelPath)
	
	// Test alias resolution
	componentsAlias := resolver.AliasPath(componentsDir)
	assert.Equal(t, "@", componentsAlias.ShortPath)
	assert.Contains(t, componentsAlias.FullPath, "@/components")
}

// TestResolverWithInvalidPaths tests resolver behavior with invalid or non-existent paths
func TestResolverWithInvalidPaths(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test with non-existent tsconfig
	nonExistentConfig := "non-existent-tsconfig.json"
	resolver := New(tmpDir, &nonExistentConfig)
	
	// Should handle gracefully and return original path
	testPath := filepath.Join(tmpDir, "src", "components")
	alias := resolver.AliasPath(testPath)
	assert.Equal(t, testPath, alias.ShortPath)
	assert.Equal(t, testPath, alias.FullPath)
	
	// Test with invalid JSON in tsconfig
	invalidTsconfigContent := `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"
    }
  }
}` // Missing closing bracket
	
	invalidTsconfigPath := filepath.Join(tmpDir, "invalid-tsconfig.json")
	err := os.WriteFile(invalidTsconfigPath, []byte(invalidTsconfigContent), 0644)
	assert.NoError(t, err)
	
	invalidTsconfigRelPath := "invalid-tsconfig.json"
	invalidResolver := New(tmpDir, &invalidTsconfigRelPath)
	
	// Should handle invalid JSON gracefully
	testPath2 := filepath.Join(tmpDir, "src", "components")
	alias2 := invalidResolver.AliasPath(testPath2)
	assert.Equal(t, testPath2, alias2.ShortPath)
	assert.Equal(t, testPath2, alias2.FullPath)
}

// TestResolverWithComplexPaths tests resolver with complex path patterns
func TestResolverWithComplexPaths(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create complex directory structure
	dirs := []string{
		filepath.Join("src", "components", "ui"),
		filepath.Join("src", "components", "layout"),
		filepath.Join("src", "utils", "helpers"),
		filepath.Join("lib", "external"),
	}
	
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		assert.NoError(t, err)
	}
	
	// Create tsconfig.json with complex path mappings
	tsconfigContent := `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@ui/*": ["src/components/ui/*"],
      "@layout/*": ["src/components/layout/*"],
      "@utils/*": ["src/utils/*"],
      "@lib/*": ["lib/*"]
    }
  }
}`
	
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err := os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Test resolver
	tsconfigRelPath := "tsconfig.json"
	resolver := New(tmpDir, &tsconfigRelPath)
	
	// Test different path resolutions
	testCases := []struct {
		path         string
		expectedAlias string
	}{
		{filepath.Join(tmpDir, "src/components/ui"), "@ui"},
		{filepath.Join(tmpDir, "src/components/layout"), "@layout"},
		{filepath.Join(tmpDir, "src/utils/helpers"), "@utils"},
		{filepath.Join(tmpDir, "lib/external"), "@lib"},
	}
	
	for _, tc := range testCases {
		alias := resolver.AliasPath(tc.path)
		assert.Equal(t, tc.expectedAlias, alias.ShortPath, "Failed for path: %s", tc.path)
	}
}

// TestResolverPathNormalization tests that resolver normalizes paths correctly
func TestResolverPathNormalization(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create tsconfig.json
	tsconfigContent := `{
  "compilerOptions": {
    "baseUrl": "./src",
    "paths": {
      "@/*": ["./*"],
      "@components/*": ["./components/*"]
    }
  }
}`
	
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err := os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Create directory structure
	srcDir := filepath.Join(tmpDir, "src")
	componentsDir := filepath.Join(srcDir, "components")
	err = os.MkdirAll(componentsDir, 0755)
	assert.NoError(t, err)
	
	// Test resolver
	tsconfigRelPath := "tsconfig.json"
	resolver := New(tmpDir, &tsconfigRelPath)
	
	// Test with different path formats
	testPaths := []string{
		filepath.Clean(componentsDir),
		filepath.ToSlash(componentsDir),
	}
	
	for _, testPath := range testPaths {
		alias := resolver.AliasPath(testPath)
		
		// Should resolve to the same alias regardless of path format
		assert.Equal(t, "@components", alias.ShortPath)
		assert.Contains(t, alias.FullPath, "@components")
		
		// FullPath should be normalized
		normalizedFullPath := filepath.ToSlash(alias.FullPath)
		assert.Equal(t, "@components", normalizedFullPath)
	}
}

// TestResolverWithSpecialCharacters tests resolver with special characters in paths
func TestResolverWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create directories with special characters
	specialDirs := []string{
		"components-with-dash",
		"components_with_underscore",
		"components.with.dots",
	}
	
	if runtime.GOOS != "windows" {
		// Unix allows more special characters
		specialDirs = append(specialDirs, "components with spaces")
	}
	
	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	assert.NoError(t, err)
	
	for _, dirName := range specialDirs {
		specialDir := filepath.Join(srcDir, dirName)
		err = os.MkdirAll(specialDir, 0755)
		assert.NoError(t, err)
	}
	
	// Create tsconfig.json with paths for special directories
	pathsConfig := make(map[string]interface{})
	for i, dirName := range specialDirs {
		aliasName := "@special" + string(rune('A'+i))
		pathsConfig[aliasName+"/*"] = []string{"src/" + dirName + "/*"}
	}
	
	tsconfigContent := `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {`
	
	for alias, paths := range pathsConfig {
		pathSlice := paths.([]string)
		tsconfigContent += `
      "` + alias + `": ["` + pathSlice[0] + `"],`
	}
	
	// Remove trailing comma and close JSON
	tsconfigContent = strings.TrimSuffix(tsconfigContent, ",")
	tsconfigContent += `
    }
  }
}`
	
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err = os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Test resolver
	tsconfigRelPath := "tsconfig.json"
	resolver := New(tmpDir, &tsconfigRelPath)
	
	// Test alias resolution for directories with special characters
	for i, dirName := range specialDirs {
		specialDir := filepath.Join(srcDir, dirName)
		alias := resolver.AliasPath(specialDir)
		
		expectedAlias := "@special" + string(rune('A'+i))
		assert.Equal(t, expectedAlias, alias.ShortPath, "Failed for directory: %s", dirName)
	}
} 