package cmd

import (
	"github.com/nergie/no-barrel-file/internal/ignorer"
	"github.com/nergie/no-barrel-file/internal/parser"

	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count barrel files in the root path",
	Run: func(cmd *cobra.Command, args []string) {
		config := NewRootConfig(cmd)
		countBarrelFiles(cmd, config)
	},
}

func countBarrelFiles(cmd *cobra.Command, config RootConfig) {
	ignorer := ignorer.New(config.rootPath, config.ignorePaths, config.gitIgnorePath)
	parser := parser.New(config.rootPath, ignorer, config.extensions)
	barrelFiles := parser.BarrelFilePaths()
	cmd.Println(len(barrelFiles))
}
