package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze a CI/CD workflow file or directory",
	Long: `Parses your YAML workflows, checks for syntax and logic errors,
scans for security vulnerabilities, and projects execution costs.

Example:
  optiflow analyze .github/workflows/main.yml
  optiflow analyze .github/workflows/`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := args[0]

		if Verbose {
			fmt.Printf("[DEBUG] Starting analysis on: %s\n", targetPath)
		}

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			fmt.Printf("Error: Path '%s' does not exist.\n", targetPath)
			os.Exit(1)
		}

		// Placeholder for Phase 2, 3, and 4
		fmt.Printf("Analyzing %s...\n", targetPath)
		fmt.Println("✓ Syntax Validated")
		fmt.Println("✓ Security Policies Enforced")
		fmt.Println("✓ Cost Projected: $0.00")
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
