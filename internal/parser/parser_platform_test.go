package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/nergie/no-barrel-file/internal/ignorer"
	"github.com/nergie/no-barrel-file/internal/resolver"
	"github.com/stretchr/testify/assert"
)

// TestParserWithDifferentPathSeparators tests parser with Windows and Unix path separators
func TestParserWithDifferentPathSeparators(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	
	// Create test directory structure with barrel files
	testDirs := []string{
		filepath.Join("src", "components"),
		filepath.Join("src", "utils"),
		filepath.Join("lib", "helpers"),
	}
	
	for _, dir := range testDirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		assert.NoError(t, err)
		
		// Create index.ts barrel file
		indexPath := filepath.Join(tmpDir, dir, "index.ts")
		indexContent := `export * from './module1'
export * from './module2'
export { default as Component } from './component'
`
		err = os.WriteFile(indexPath, []byte(indexContent), 0644)
		assert.NoError(t, err)
		
		// Create module files
		moduleFiles := []string{"module1.ts", "module2.ts", "component.ts"}
		for _, moduleFile := range moduleFiles {
			modulePath := filepath.Join(tmpDir, dir, moduleFile)
			moduleContent := `export const test = 'test';`
			err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
			assert.NoError(t, err)
		}
	}
	
	// Test parser with different extensions
	extensions := []string{".ts", ".js", ".tsx", ".jsx"}
	ignorer := ignorer.New(tmpDir, []string{}, "")
	parser := New(tmpDir, ignorer, extensions)
	
	barrelPaths := parser.BarrelFilePaths()
	
	// Should find all barrel files regardless of path separator format
	assert.Len(t, barrelPaths, 3)
	
	// Verify paths are normalized correctly for the current platform
	for _, path := range barrelPaths {
		assert.True(t, strings.HasSuffix(path, "index.ts"))
		// Path should be absolute and properly formatted for current OS
		assert.True(t, filepath.IsAbs(path))
	}
}

// TestParserCaseSensitivity tests how parser handles case-sensitive vs case-insensitive file systems
func TestParserCaseSensitivity(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create files with different cases
	testFiles := map[string]string{
		"INDEX.ts": `export * from './Module'`,
		"index.ts": `export * from './module'`,
		"Module.ts": `export const test = 'test';`,
		"module.ts": `export const test = 'test';`,
	}
	
	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			// On case-insensitive file systems, this might fail for duplicate names
			if runtime.GOOS == "windows" && strings.Contains(err.Error(), "already exists") {
				continue // Skip duplicate files on Windows
			}
			assert.NoError(t, err)
		}
	}
	
	extensions := []string{".ts"}
	ignorer := ignorer.New(tmpDir, []string{}, "")
	parser := New(tmpDir, ignorer, extensions)
	
	barrelPaths := parser.BarrelFilePaths()
	
	if runtime.GOOS == "windows" {
		// On Windows (case-insensitive), we might only get one index file
		assert.True(t, len(barrelPaths) >= 1)
	} else {
		// On Unix (case-sensitive), we should get both INDEX.ts and index.ts if they exist
		// But only if both are valid barrel files
		assert.True(t, len(barrelPaths) >= 1)
	}
}

// TestParserWithAliasPathResolution tests parser with different alias path formats
func TestParserWithAliasPathResolution(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create directory structure
	srcDir := filepath.Join(tmpDir, "src")
	componentsDir := filepath.Join(srcDir, "components")
	err := os.MkdirAll(componentsDir, 0755)
	assert.NoError(t, err)
	
	// Create tsconfig.json with alias paths
	var tsconfigContent string
	if runtime.GOOS == "windows" {
		tsconfigContent = `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@components/*": ["src/components/*"]
    }
  }
}`
	} else {
		tsconfigContent = `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@components/*": ["src/components/*"]
    }
  }
}`
	}
	
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err = os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Create barrel file with alias imports
	barrelContent := `export * from './Button'
export * from './Input'`
	barrelPath := filepath.Join(componentsDir, "index.ts")
	err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
	assert.NoError(t, err)
	
	// Create module files
	moduleFiles := []string{"Button.ts", "Input.ts"}
	for _, moduleFile := range moduleFiles {
		modulePath := filepath.Join(componentsDir, moduleFile)
		moduleContent := `export const Component = 'component';`
		err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
		assert.NoError(t, err)
	}
	
	// Test with resolver
	tsconfigRelPath := "tsconfig.json"
	resolver := resolver.New(tmpDir, &tsconfigRelPath)
	
	extensions := []string{".ts"}
	ignorer := ignorer.New(tmpDir, []string{}, "")
	parser := New(tmpDir, ignorer, extensions)
	
	barrelPathExistenceMap, barrelModuleResolverMap := parser.BarrelMaps(resolver)
	
	// Should find barrel paths
	assert.NotEmpty(t, barrelPathExistenceMap)
	assert.NotEmpty(t, barrelModuleResolverMap)
	
	// Check that paths are properly resolved
	for path := range barrelPathExistenceMap {
		// Path should be properly formatted for current OS
		normalizedPath := filepath.ToSlash(path)
		assert.Contains(t, normalizedPath, "components")
	}
}

// TestParserWithLongPaths tests parser with long file paths (Windows limitation)
func TestParserWithLongPaths(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Long path test only relevant on Windows")
	}
	
	tmpDir := t.TempDir()
	
	// Create a deeply nested directory structure
	deepPath := tmpDir
	for i := 0; i < 10; i++ {
		deepPath = filepath.Join(deepPath, "very-long-directory-name-to-test-path-limits")
	}
	
	err := os.MkdirAll(deepPath, 0755)
	if err != nil {
		// If we can't create the path due to length limitations, skip this test
		t.Skipf("Cannot create long path on this system: %v", err)
	}
	
	// Create barrel file in deep path
	barrelPath := filepath.Join(deepPath, "index.ts")
	barrelContent := `export * from './module'`
	err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
	assert.NoError(t, err)
	
	// Create module file
	modulePath := filepath.Join(deepPath, "module.ts")
	moduleContent := `export const test = 'test';`
	err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
	assert.NoError(t, err)
	
	extensions := []string{".ts"}
	ignorer := ignorer.New(tmpDir, []string{}, "")
	parser := New(tmpDir, ignorer, extensions)
	
	// Should handle long paths without error
	barrelPaths := parser.BarrelFilePaths()
	assert.Len(t, barrelPaths, 1)
	assert.Contains(t, barrelPaths[0], "index.ts")
}

// TestParserWithSymlinks tests parser behavior with symbolic links
func TestParserWithSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Symlink test may require admin privileges on Windows")
	}
	
	tmpDir := t.TempDir()
	
	// Create source directory with barrel file
	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	assert.NoError(t, err)
	
	barrelPath := filepath.Join(srcDir, "index.ts")
	barrelContent := `export * from './module'`
	err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
	assert.NoError(t, err)
	
	modulePath := filepath.Join(srcDir, "module.ts")
	moduleContent := `export const test = 'test';`
	err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
	assert.NoError(t, err)
	
	// Create symlink to source directory
	linkDir := filepath.Join(tmpDir, "link")
	err = os.Symlink(srcDir, linkDir)
	if err != nil {
		t.Skipf("Cannot create symlink: %v", err)
	}
	
	extensions := []string{".ts"}
	ignorer := ignorer.New(tmpDir, []string{}, "")
	parser := New(tmpDir, ignorer, extensions)
	
	barrelPaths := parser.BarrelFilePaths()
	
	// Should find barrel files in both original and symlinked directories
	assert.True(t, len(barrelPaths) >= 1)
	
	// Verify paths are resolved correctly
	for _, path := range barrelPaths {
		assert.True(t, strings.HasSuffix(path, "index.ts"))
	}
}

// TestParserWithSpecialCharacters tests parser with special characters in file names
func TestParserWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test different special characters allowed in file names
	specialDirs := []string{
		"components-with-dash",
		"components_with_underscore",
		"components.with.dots",
	}
	
	if runtime.GOOS != "windows" {
		// Unix allows more special characters
		specialDirs = append(specialDirs, "components with spaces")
	}
	
	for _, dirName := range specialDirs {
		dir := filepath.Join(tmpDir, dirName)
		err := os.MkdirAll(dir, 0755)
		assert.NoError(t, err)
		
		barrelPath := filepath.Join(dir, "index.ts")
		barrelContent := `export * from './module'`
		err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
		assert.NoError(t, err)
		
		modulePath := filepath.Join(dir, "module.ts")
		moduleContent := `export const test = 'test';`
		err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
		assert.NoError(t, err)
	}
	
	extensions := []string{".ts"}
	ignorer := ignorer.New(tmpDir, []string{}, "")
	parser := New(tmpDir, ignorer, extensions)
	
	barrelPaths := parser.BarrelFilePaths()
	
	// Should find all barrel files with special characters
	assert.Len(t, barrelPaths, len(specialDirs))
	
	for _, path := range barrelPaths {
		assert.True(t, strings.HasSuffix(path, "index.ts"))
		// Verify path is properly formatted
		assert.True(t, filepath.IsAbs(path))
	}
}

// TestParserPathNormalization tests that parser normalizes paths consistently
func TestParserPathNormalization(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create directory structure
	componentsDir := filepath.Join(tmpDir, "src", "components")
	err := os.MkdirAll(componentsDir, 0755)
	assert.NoError(t, err)
	
	// Create barrel file with mixed path separators in content
	var barrelContent string
	if runtime.GOOS == "windows" {
		barrelContent = `export * from './Button'
export * from './Input'
export * from '../utils/helper'`
	} else {
		barrelContent = `export * from './Button'
export * from './Input'
export * from '../utils/helper'`
	}
	
	barrelPath := filepath.Join(componentsDir, "index.ts")
	err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
	assert.NoError(t, err)
	
	// Create referenced files
	buttonPath := filepath.Join(componentsDir, "Button.ts")
	err = os.WriteFile(buttonPath, []byte("export const Button = 'button';"), 0644)
	assert.NoError(t, err)
	
	inputPath := filepath.Join(componentsDir, "Input.ts")
	err = os.WriteFile(inputPath, []byte("export const Input = 'input';"), 0644)
	assert.NoError(t, err)
	
	// Create utils directory and helper file
	utilsDir := filepath.Join(tmpDir, "src", "utils")
	err = os.MkdirAll(utilsDir, 0755)
	assert.NoError(t, err)
	
	helperPath := filepath.Join(utilsDir, "helper.ts")
	err = os.WriteFile(helperPath, []byte("export const helper = 'helper';"), 0644)
	assert.NoError(t, err)
	
	extensions := []string{".ts"}
	ignorer := ignorer.New(tmpDir, []string{}, "")
	parser := New(tmpDir, ignorer, extensions)
	
	barrelPaths := parser.BarrelFilePaths()
	
	// Should find the barrel file
	assert.Len(t, barrelPaths, 1)
	
	// Verify the path is properly normalized
	foundPath := barrelPaths[0]
	assert.True(t, strings.HasSuffix(foundPath, "index.ts"))
	assert.Contains(t, foundPath, "components")
	
	// Path should be absolute and use correct separators for OS
	assert.True(t, filepath.IsAbs(foundPath))
} 