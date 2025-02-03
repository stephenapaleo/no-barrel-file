package cmd_flag

import (
	"strings"

	"github.com/spf13/cobra"
)

func AddRootPath(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		"root-path", "r", "", "Root path of the targeted project.")
	cmd.MarkPersistentFlagRequired("root-path")

}

func RootPath(cmd *cobra.Command) string {
	return cmd.Flags().Lookup("root-path").Value.String()
}

func AddExtensions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("extensions", "e", ".ts,.js,.tsx,.jsx", "Comma-separated list of file extensions to process.")
}

func Extensions(cmd *cobra.Command) []string {
	extensionsString := cmd.Flags().Lookup("extensions").Value.String()
	extensionsString = strings.TrimSpace(extensionsString)
	if extensionsString == "" {
		return []string{}
	}
	return strings.Split(extensionsString, ",")
}

func AddGitIgnorePath(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		"gitignore-path", "g", ".gitignore", "Relative path to `.gitignore` file to apply ignore rules.")
}

func GitIgnorePath(cmd *cobra.Command) string {
	return cmd.Flags().Lookup("gitignore-path").Value.String()
}

func AddIgnorePaths(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		"ignore-paths", "i", "", "Comma-separated list of directories or files to ignore.")
}

func IgnorePaths(cmd *cobra.Command) []string {
	pathsString := cmd.Flags().Lookup("ignore-paths").Value.String()
	if pathsString == "" {
		return []string{}
	}
	return strings.Split(pathsString, ",")
}

func AddTargetPath(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(
		"target-path", "t", ".", "Relative path where imports should be replaced.")
}

func TargetPath(cmd *cobra.Command) string {
	return cmd.Flags().Lookup("target-path").Value.String()
}

func AddAliasConfigPath(cmd *cobra.Command) {
	cmd.Flags().StringP(
		"alias-config-path", "a", "", "Relative path to 'tsconfig.json' or 'jsconfig.json' for alias resolution. Only JSON files are supported.")
}

func AliasConfigPath(cmd *cobra.Command) string {
	return cmd.Flags().Lookup("alias-config-path").Value.String()
}

func AddBarrelPath(cmd *cobra.Command) {
	cmd.Flags().StringP(
		"barrel-path", "b", ".", "Relative path of a barrel file import to replaced.")
}

func BarrelPath(cmd *cobra.Command) string {
	return cmd.Flags().Lookup("barrel-path").Value.String()
}

func AddVerbose(cmd *cobra.Command) {
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

func Verbose(cmd *cobra.Command) bool {
	isVerbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return false
	}
	return isVerbose
}
