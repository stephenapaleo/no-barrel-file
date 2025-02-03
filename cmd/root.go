package cmd

import (
	"fmt"
	"os"

	"github.com/nergie/no-barrel-file/internal/cmd_flag"

	"github.com/spf13/cobra"
)

type RootConfig struct {
	gitIgnorePath string
	ignorePaths   []string
	rootPath      string
	extensions    []string
}

func NewRootConfig(cmd *cobra.Command) RootConfig {
	return RootConfig{
		gitIgnorePath: cmd_flag.GitIgnorePath(cmd),
		ignorePaths:   cmd_flag.IgnorePaths(cmd),
		rootPath:      cmd_flag.RootPath(cmd),
		extensions:    cmd_flag.Extensions(cmd),
	}
}

var (
	rootCmd = &cobra.Command{
		Use:   "barrel-file",
		Short: "A CLI tool for managing barrel files",
		Long:  `no-barrel-file is a CLI tool to replace barrel imports, count, and display barrel files in folders.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cmd_flag.AddIgnorePaths(rootCmd)
	cmd_flag.AddGitIgnorePath(rootCmd)
	cmd_flag.AddExtensions(rootCmd)
	cmd_flag.AddRootPath(rootCmd)

	rootCmd.AddCommand(countCmd)
	rootCmd.AddCommand(displayCmd)
	rootCmd.AddCommand(replaceCmd)
}
