package cmd

import (
	"path/filepath"

	"github.com/nergie/no-barrel-file/internal/ignorer"
	"github.com/nergie/no-barrel-file/internal/parser"

	"github.com/spf13/cobra"
)

var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "Display barrel files in the root path",
	Run: func(cmd *cobra.Command, args []string) {
		config := NewRootConfig(cmd)
		displayBarrelFiles(cmd, config)
	},
}

func displayBarrelFiles(cmd *cobra.Command, config RootConfig) {
	ignorer := ignorer.New(config.rootPath, config.ignorePaths, config.gitIgnorePath)
	parser := parser.New(config.rootPath, ignorer, config.extensions)
	barrelPaths := parser.BarrelFilePaths()
	cmd.Printf("%d barrel files found\n", len(barrelPaths))
	for _, fullPath := range barrelPaths {
		relativePath, err := filepath.Rel(config.rootPath, fullPath)
		if err == nil {
			cmd.Println(relativePath)
		}
	}
}
