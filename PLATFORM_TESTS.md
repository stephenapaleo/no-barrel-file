# Platform-Specific Unit Tests for Windows and Linux

This document describes the comprehensive unit tests that have been created to ensure the `no-barrel-file` CLI tool works correctly on both Windows and Linux file systems.

## Overview

The tool processes JavaScript/TypeScript barrel files and replaces barrel imports with direct imports. Since Windows and Linux file systems behave differently, we need tests to ensure cross-platform compatibility.

## Key File System Differences Tested

### 1. Path Separators
- **Windows**: Uses backslashes (`\`) as path separators
- **Linux**: Uses forward slashes (`/`) as path separators
- **Tests**: Verify that the tool correctly handles both path separator formats

### 2. Case Sensitivity
- **Windows**: File system is case-insensitive (`index.ts` and `INDEX.ts` are the same file)
- **Linux**: File system is case-sensitive (`index.ts` and `INDEX.ts` are different files)
- **Tests**: Verify that the tool behaves correctly on both types of file systems

### 3. Path Normalization
- **Windows**: Uses drive letters (`C:\`) and has specific path length limitations
- **Linux**: Uses single root (`/`) and allows longer paths
- **Tests**: Verify that paths are normalized correctly for each platform

### 4. Special Characters
- **Windows**: More restrictive about special characters in file names
- **Linux**: Allows more special characters including spaces
- **Tests**: Verify that the tool handles different special characters appropriately

## Test Files Created

### 1. `internal/tests/platform_test.go`
Basic platform-specific tests for fundamental file system operations:
- `TestPathNormalization`: Tests path normalization across platforms
- `TestPathJoin`: Tests cross-platform path joining
- `TestPathMatching`: Tests case sensitivity differences
- `TestAbsolutePaths`: Tests absolute path handling
- `TestRelativePaths`: Tests relative path calculation
- `TestPathSeparatorConversion`: Tests path separator conversion
- `TestFileSystemOperations`: Tests file system operations across platforms
- `TestPathPattern`: Tests path pattern matching

### 2. `internal/parser/parser_platform_test.go`
Platform-specific tests for the parser component:
- `TestParserWithDifferentPathSeparators`: Tests parser with Windows and Unix path separators
- `TestParserCaseSensitivity`: Tests how parser handles case-sensitive vs case-insensitive file systems
- `TestParserWithAliasPathResolution`: Tests parser with different alias path formats
- `TestParserWithLongPaths`: Tests parser with long file paths (Windows limitation)
- `TestParserWithSymlinks`: Tests parser behavior with symbolic links (Unix specific)
- `TestParserWithSpecialCharacters`: Tests parser with special characters in file names
- `TestParserPathNormalization`: Tests that parser normalizes paths consistently

### 3. `internal/resolver/resolver_platform_test.go`
Platform-specific tests for the resolver component:
- `TestResolverWithWindowsPaths`: Tests resolver with Windows-style paths
- `TestResolverWithRelativePaths`: Tests resolver with different relative path formats
- `TestResolverCaseSensitivity`: Tests resolver with case-sensitive vs case-insensitive paths
- `TestResolverWithAbsolutePaths`: Tests resolver with absolute paths in tsconfig
- `TestResolverWithInvalidPaths`: Tests resolver behavior with invalid or non-existent paths
- `TestResolverWithComplexPaths`: Tests resolver with complex path patterns
- `TestResolverPathNormalization`: Tests that resolver normalizes paths correctly
- `TestResolverWithSpecialCharacters`: Tests resolver with special characters in paths

### 4. `cmd/platform_test.go`
Platform-specific tests for CLI commands:
- `TestCountCommandWithDifferentPathFormats`: Tests count command with different path formats
- `TestDisplayCommandCaseSensitivity`: Tests display command with case sensitivity differences
- `TestReplaceCommandWithDifferentSeparators`: Tests replace command with different path separators
- `TestCommandsWithLongPaths`: Tests commands with long file paths (Windows specific)
- `TestCommandsWithSpecialCharacters`: Tests commands with special characters in paths
- `TestCommandsWithSymlinks`: Tests commands with symbolic links (Unix specific)
- `TestCommandsWithIgnorePatterns`: Tests commands with gitignore patterns and different path formats
- `TestCommandsWithMixedCaseExtensions`: Tests commands with mixed case file extensions

### 5. `internal/ignorer/ignorer_platform_test.go`
Platform-specific tests for the ignorer component:
- `TestIgnorerWithDifferentPathSeparators`: Tests ignorer with Windows and Unix path separators
- `TestIgnorerCaseSensitivity`: Tests ignorer with case-sensitive vs case-insensitive patterns
- `TestIgnorerWithAbsolutePaths`: Tests ignorer with absolute vs relative paths
- `TestIgnorerWithCustomIgnorePaths`: Tests ignorer with custom ignore paths using different separators
- `TestIgnorerWithComplexPatterns`: Tests ignorer with complex gitignore patterns
- `TestIgnorerWithSpecialCharacters`: Tests ignorer with special characters in paths
- `TestIgnorerWithSymlinks`: Tests ignorer behavior with symbolic links

## Running the Tests

### Run All Platform Tests
```bash
# Run all platform-specific tests
go test ./internal/tests/... -v
go test ./internal/parser/... -v
go test ./internal/resolver/... -v
go test ./cmd/... -v
go test ./internal/ignorer/... -v
```

### Run Tests on Specific Platform
```bash
# Run tests only on Windows
go test ./... -v -run="Windows"

# Run tests only on Unix/Linux
go test ./... -v -run="Unix"

# Run symlink tests (Unix only)
go test ./... -v -run="Symlink"

# Run long path tests (Windows only)
go test ./... -v -run="LongPath"
```

## Test Scenarios Covered

### Path Handling
- ✅ Forward slash paths (`src/components/index.ts`)
- ✅ Backslash paths (`src\components\index.ts`)
- ✅ Mixed separator paths (`src/components\index.ts`)
- ✅ Absolute paths with drive letters (Windows)
- ✅ Absolute paths with root slash (Unix)
- ✅ Relative paths (`./components`, `../utils`)
- ✅ Path normalization and cleaning

### Case Sensitivity
- ✅ Exact case matches
- ✅ Different case files on case-sensitive systems
- ✅ Case-insensitive matching on Windows
- ✅ Case-sensitive matching on Linux
- ✅ Mixed case file extensions

### File System Features
- ✅ Long path handling (Windows specific)
- ✅ Symbolic link handling (Unix specific)
- ✅ Special characters in file names
- ✅ Directory traversal with different separators
- ✅ Gitignore pattern matching

### TypeScript Configuration
- ✅ Alias path resolution with different formats
- ✅ Relative and absolute baseUrl configurations
- ✅ Complex path mapping patterns
- ✅ Invalid or missing tsconfig handling

### CLI Commands
- ✅ Count command with different path formats
- ✅ Display command with case sensitivity
- ✅ Replace command with path separators
- ✅ Ignore patterns with different formats

## Platform-Specific Behaviors Tested

### Windows-Specific
- Case-insensitive file system behavior
- Backslash path separators in gitignore
- Drive letter absolute paths
- Long path limitations (MAX_PATH)
- Restricted special characters in file names

### Linux-Specific
- Case-sensitive file system behavior
- Forward slash path separators
- Root-based absolute paths
- Symbolic link resolution
- Extended special characters in file names

## Benefits

These comprehensive tests ensure that:

1. **Cross-Platform Compatibility**: The tool works correctly on both Windows and Linux
2. **Path Handling**: All path formats are properly normalized and processed
3. **File System Awareness**: The tool adapts to different file system behaviors
4. **Error Handling**: Invalid paths and configurations are handled gracefully
5. **Performance**: Long paths and complex directory structures are handled efficiently
6. **Reliability**: Edge cases specific to each platform are properly tested

## Maintenance

When adding new features or modifying existing functionality:

1. Consider how the changes might behave differently on Windows vs Linux
2. Add corresponding platform-specific tests
3. Run tests on both platforms before merging
4. Update this documentation with new test scenarios

The platform-specific tests provide confidence that the tool will work reliably for users on both Windows and Linux systems, handling the various file system differences transparently. 