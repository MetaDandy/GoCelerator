package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show goce and Go versions",
	Long:  "Prints the current version of goce and the Go runtime version.",
	Run: func(cmd *cobra.Command, args []string) {
		info, ok := debug.ReadBuildInfo()
		var ver string
		if ok && info.Main.Version != "" {
			ver = info.Main.Version
		} else {
			ver = version
		}
		fmt.Printf("goce version: %s\n", ver)
		fmt.Printf("Go version:   %s\n", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
