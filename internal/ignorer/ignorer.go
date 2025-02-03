package ignorer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

type Ignorer struct {
	gitIgnore       *ignore.GitIgnore
	manualIgnoreMap map[string]struct{ IsDir bool }
	rootPath        string
}

func New(rootPath string, ignorePaths []string, gitIgnorePath string) Ignorer {
	gitIgnore := getGitignore(rootPath, gitIgnorePath)
	manualIgnoreMap := getManualIgnore(rootPath, ignorePaths)
	return Ignorer{
		gitIgnore:       gitIgnore,
		manualIgnoreMap: manualIgnoreMap,
		rootPath:        rootPath,
	}
}

func (ignorer *Ignorer) IgnorePath(path string) bool {
	relativePath, _ := filepath.Rel(ignorer.rootPath, path)
	if ignorer.gitIgnore != nil && ignorer.gitIgnore.MatchesPath(relativePath) {
		return true
	}
	if _, ignored := ignorer.manualIgnoreMap[path]; ignored {
		return true
	}
	for dir, isDir := range ignorer.manualIgnoreMap {
		if isDir.IsDir && strings.HasPrefix(path, dir) {
			return true
		}
	}
	return false
}

func getGitignore(rootPath string, gitIgnorePath string) *ignore.GitIgnore {
	gitIgnoreFullPath := filepath.Join(rootPath, gitIgnorePath)
	file, err := os.Open(gitIgnoreFullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find gitignore file: %v\n", err)
		fmt.Fprintln(os.Stderr, "Ignoring gitignore file")
		return nil
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .gitignore file: %v\n", err)
		return nil
	}

	return ignore.CompileIgnoreLines(patterns...)
}

func getManualIgnore(rootPath string, ignorePaths []string) map[string]struct{ IsDir bool } {
	manualIgnoreMap := make(map[string]struct{ IsDir bool })
	for _, path := range ignorePaths {
		fullPath := filepath.Join(rootPath, path)
		if fileInfo, err := os.Stat(fullPath); err == nil && fileInfo.IsDir() {
			manualIgnoreMap[fullPath] = struct{ IsDir bool }{IsDir: true}
		} else {
			manualIgnoreMap[fullPath] = struct{ IsDir bool }{IsDir: false}
		}
	}
	return manualIgnoreMap
}
