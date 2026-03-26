package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GoReleaser at build time via ldflags
var AppVersion = "v0.1.0-dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of OptiFlow-CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("OptiFlow-CLI %s\n", AppVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
