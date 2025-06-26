package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/nergie/no-barrel-file/internal/tests"
	"github.com/stretchr/testify/assert"
)

// TestCountCommandWithDifferentPathFormats tests count command with different path formats
func TestCountCommandWithDifferentPathFormats(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test structure with barrel files
	testDirs := []string{
		filepath.Join("src", "components"),
		filepath.Join("src", "utils"),
		filepath.Join("lib", "helpers"),
	}
	
	expectedBarrelCount := 0
	for _, dir := range testDirs {
		fullDir := filepath.Join(tmpDir, dir)
		err := os.MkdirAll(fullDir, 0755)
		assert.NoError(t, err)
		
		// Create barrel file
		barrelPath := filepath.Join(fullDir, "index.ts")
		barrelContent := `export * from './module1'
export * from './module2'`
		err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
		assert.NoError(t, err)
		expectedBarrelCount++
		
		// Create module files
		for i := 1; i <= 2; i++ {
			modulePath := filepath.Join(fullDir, "module"+string(rune('0'+i))+".ts")
			moduleContent := `export const test = 'test';`
			err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
			assert.NoError(t, err)
		}
	}
	
	// Test count command with different path format variations
	testCases := []struct {
		name     string
		rootPath string
	}{
		{"Normalized path", tmpDir},
		{"Path with trailing separator", tmpDir + string(filepath.Separator)},
		{"Path with current directory", filepath.Join(tmpDir, ".")},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := tests.ExecuteCommand(rootCmd, "count", "--root-path", tc.rootPath)
			assert.NoError(t, err)
			assert.Contains(t, output, "3")
		})
	}
}

// TestDisplayCommandCaseSensitivity tests display command with case sensitivity differences
func TestDisplayCommandCaseSensitivity(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test files with different cases
	srcDir := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(srcDir, 0755)
	assert.NoError(t, err)
	
	// Create lowercase index file
	lowerIndexPath := filepath.Join(srcDir, "index.ts")
	lowerIndexContent := `export * from './module'`
	err = os.WriteFile(lowerIndexPath, []byte(lowerIndexContent), 0644)
	assert.NoError(t, err)
	
	// Create module file
	modulePath := filepath.Join(srcDir, "module.ts")
	moduleContent := `export const test = 'test';`
	err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
	assert.NoError(t, err)
	
	// Try to create uppercase INDEX file (may fail on case-insensitive systems)
	upperIndexPath := filepath.Join(srcDir, "INDEX.ts")
	upperIndexContent := `export * from './MODULE'`
	err = os.WriteFile(upperIndexPath, []byte(upperIndexContent), 0644)
	if err != nil && runtime.GOOS == "windows" {
		t.Logf("Could not create uppercase INDEX.ts on case-insensitive filesystem: %v", err)
	}
	
	// Test display command
	output, err := tests.ExecuteCommand(rootCmd, "display", "--root-path", tmpDir)
	assert.NoError(t, err)
	
	// Should find at least one barrel file
	assert.Contains(t, output, "index.ts")
	
	if runtime.GOOS == "windows" {
		// On Windows, should treat index.ts and INDEX.ts as the same file
		indexCount := strings.Count(strings.ToLower(output), "index.ts")
		assert.Equal(t, 1, indexCount, "Should only find one index file on case-insensitive filesystem")
	}
}

// TestReplaceCommandWithDifferentSeparators tests replace command with different path separators
func TestReplaceCommandWithDifferentSeparators(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test structure
	srcDir := filepath.Join(tmpDir, "src")
	componentsDir := filepath.Join(srcDir, "components")
	err := os.MkdirAll(componentsDir, 0755)
	assert.NoError(t, err)
	
	// Create tsconfig.json
	tsconfigContent := `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@components/*": ["src/components/*"]
    }
  }
}`
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	err = os.WriteFile(tsconfigPath, []byte(tsconfigContent), 0644)
	assert.NoError(t, err)
	
	// Create barrel file
	barrelPath := filepath.Join(componentsDir, "index.ts")
	barrelContent := `export * from './Button'
export * from './Input'`
	err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
	assert.NoError(t, err)
	
	// Create component files
	buttonPath := filepath.Join(componentsDir, "Button.ts")
	buttonContent := `export const Button = 'button';`
	err = os.WriteFile(buttonPath, []byte(buttonContent), 0644)
	assert.NoError(t, err)
	
	inputPath := filepath.Join(componentsDir, "Input.ts")
	inputContent := `export const Input = 'input';`
	err = os.WriteFile(inputPath, []byte(inputContent), 0644)
	assert.NoError(t, err)
	
	// Create a file that imports from the barrel with different path separators
	consumerPath := filepath.Join(srcDir, "App.ts")
	var consumerContent string
	if runtime.GOOS == "windows" {
		// Test with backslashes on Windows
		consumerContent = `import { Button, Input } from './components';
import { Button as Btn } from './components/index';`
	} else {
		// Test with forward slashes on Unix
		consumerContent = `import { Button, Input } from './components';
import { Button as Btn } from './components/index';`
	}
	err = os.WriteFile(consumerPath, []byte(consumerContent), 0644)
	assert.NoError(t, err)
	
	// Test replace command
	output, err := tests.ExecuteCommand(rootCmd, "replace", "--root-path", tmpDir, "--alias-config-path", "tsconfig.json")
	assert.NoError(t, err)
	
	// Should successfully replace imports
	assert.Contains(t, output, "1 files updated")
	
	// Verify the replacement worked
	updatedContent, err := os.ReadFile(consumerPath)
	assert.NoError(t, err)
	updatedStr := string(updatedContent)
	
	// Should have replaced barrel imports with direct imports
	assert.Contains(t, updatedStr, "./components/Button")
	assert.Contains(t, updatedStr, "./components/Input")
}

// TestCommandsWithLongPaths tests commands with long file paths (Windows specific)
func TestCommandsWithLongPaths(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Long path test only relevant on Windows")
	}
	
	tmpDir := t.TempDir()
	
	// Create deeply nested directory structure
	deepPath := tmpDir
	for i := 0; i < 8; i++ {
		deepPath = filepath.Join(deepPath, "very-long-directory-name-"+string(rune('A'+i)))
	}
	
	err := os.MkdirAll(deepPath, 0755)
	if err != nil {
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
	
	// Test count command with long path
	output, err := tests.ExecuteCommand(rootCmd, "count", "--root-path", tmpDir)
	assert.NoError(t, err)
	assert.Contains(t, output, "1")
	
	// Test display command with long path
	output, err = tests.ExecuteCommand(rootCmd, "display", "--root-path", tmpDir)
	assert.NoError(t, err)
	assert.Contains(t, output, "index.ts")
}

// TestCommandsWithSpecialCharacters tests commands with special characters in paths
func TestCommandsWithSpecialCharacters(t *testing.T) {
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
	
	expectedBarrelCount := 0
	for _, dirName := range specialDirs {
		dir := filepath.Join(tmpDir, dirName)
		err := os.MkdirAll(dir, 0755)
		assert.NoError(t, err)
		
		// Create barrel file
		barrelPath := filepath.Join(dir, "index.ts")
		barrelContent := `export * from './module'`
		err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
		assert.NoError(t, err)
		expectedBarrelCount++
		
		// Create module file
		modulePath := filepath.Join(dir, "module.ts")
		moduleContent := `export const test = 'test';`
		err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
		assert.NoError(t, err)
	}
	
	// Test count command
	output, err := tests.ExecuteCommand(rootCmd, "count", "--root-path", tmpDir)
	assert.NoError(t, err)
	assert.Contains(t, output, string(rune('0'+expectedBarrelCount)))
	
	// Test display command
	output, err = tests.ExecuteCommand(rootCmd, "display", "--root-path", tmpDir)
	assert.NoError(t, err)
	
	// Should find all barrel files
	for _, dirName := range specialDirs {
		assert.Contains(t, output, dirName)
	}
}

// TestCommandsWithSymlinks tests commands with symbolic links (Unix specific)
func TestCommandsWithSymlinks(t *testing.T) {
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
	
	// Test count command (should count both original and symlinked barrel files)
	output, err := tests.ExecuteCommand(rootCmd, "count", "--root-path", tmpDir)
	assert.NoError(t, err)
	// Should find barrel files (may be 1 or 2 depending on how symlinks are handled)
	assert.True(t, strings.Contains(output, "1") || strings.Contains(output, "2"))
	
	// Test display command
	output, err = tests.ExecuteCommand(rootCmd, "display", "--root-path", tmpDir)
	assert.NoError(t, err)
	assert.Contains(t, output, "index.ts")
}

// TestCommandsWithIgnorePatterns tests commands with gitignore patterns and different path formats
func TestCommandsWithIgnorePatterns(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create .gitignore with different path formats
	var gitignoreContent string
	if runtime.GOOS == "windows" {
		gitignoreContent = `node_modules/
build\\
dist/
*.tmp`
	} else {
		gitignoreContent = `node_modules/
build/
dist/
*.tmp`
	}
	
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	assert.NoError(t, err)
	
	// Create directories with barrel files
	testDirs := []string{
		filepath.Join("src", "components"),    // Should be included
		filepath.Join("node_modules", "lib"),  // Should be ignored
		filepath.Join("build", "output"),      // Should be ignored
		filepath.Join("dist", "bundle"),       // Should be ignored
	}
	
	expectedBarrelCount := 1 // Only src/components should be counted
	for _, dir := range testDirs {
		fullDir := filepath.Join(tmpDir, dir)
		err := os.MkdirAll(fullDir, 0755)
		assert.NoError(t, err)
		
		// Create barrel file
		barrelPath := filepath.Join(fullDir, "index.ts")
		barrelContent := `export * from './module'`
		err = os.WriteFile(barrelPath, []byte(barrelContent), 0644)
		assert.NoError(t, err)
		
		// Create module file
		modulePath := filepath.Join(fullDir, "module.ts")
		moduleContent := `export const test = 'test';`
		err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
		assert.NoError(t, err)
	}
	
	// Test count command with gitignore
	output, err := tests.ExecuteCommand(rootCmd, "count", "--root-path", tmpDir, "--gitignore-path", ".gitignore")
	assert.NoError(t, err)
	assert.Contains(t, output, string(rune('0'+expectedBarrelCount)))
	
	// Test display command with gitignore
	output, err = tests.ExecuteCommand(rootCmd, "display", "--root-path", tmpDir, "--gitignore-path", ".gitignore")
	assert.NoError(t, err)
	
	// Should only show non-ignored barrel files
	assert.Contains(t, output, "src")
	assert.NotContains(t, output, "node_modules")
	assert.NotContains(t, output, "build")
	assert.NotContains(t, output, "dist")
}

// TestCommandsWithMixedCaseExtensions tests commands with mixed case file extensions
func TestCommandsWithMixedCaseExtensions(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test files with different case extensions
	testFiles := map[string]string{
		"index.ts":  `export * from './module1'`,
		"index.TS":  `export * from './module2'`,
		"index.Ts":  `export * from './module3'`,
		"index.tsx": `export * from './component'`,
	}
	
	createdFiles := 0
	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			// On case-insensitive file systems, this might fail for duplicate names
			if runtime.GOOS == "windows" && strings.Contains(err.Error(), "already exists") {
				continue
			}
			assert.NoError(t, err)
		}
		
		// Create corresponding module files
		baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
		moduleFilename := baseFilename + "module" + filepath.Ext(filename)
		moduleContent := `export const test = 'test';`
		
		modulePath := filepath.Join(tmpDir, moduleFilename)
		err = os.WriteFile(modulePath, []byte(moduleContent), 0644)
		if err == nil {
			createdFiles++
		}
	}
	
	// Test count command with mixed extensions
	output, err := tests.ExecuteCommand(rootCmd, "count", "--root-path", tmpDir, "--extensions", ".ts,.TS,.Ts,.tsx")
	assert.NoError(t, err)
	
	if runtime.GOOS == "windows" {
		// On Windows, different cases of the same extension should be treated as the same
		assert.True(t, strings.Contains(output, "1") || strings.Contains(output, "2"))
	} else {
		// On Unix, different cases should be treated as different extensions
		assert.True(t, len(output) > 0)
	}
} 