package cmd

import (
	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:   "optiflow",
	Short: "A blazing-fast, local CI/CD pipeline analyzer",
	Long: `OptiFlow-CLI validates your CI/CD pipelines locally.
It provides instant logic linting, security enforcement, 
and shift-left cost projections before you even commit your code.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output for debugging")
}
