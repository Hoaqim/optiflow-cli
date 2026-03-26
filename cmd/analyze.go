package cmd

import (
	"fmt"
	"os"

	"github.com/Hoaqim/optiflow-cli/pkg/finops"
	"github.com/Hoaqim/optiflow-cli/pkg/security"
	"github.com/Hoaqim/optiflow-cli/pkg/workflow"
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

		wf, err := workflow.Parse(targetPath)
		if err != nil {
			fmt.Printf("Syntax Error in %s:\n%v\n", targetPath, err)
			os.Exit(1)
		}

		fmt.Printf("Analyzing %s...\n", targetPath)
		fmt.Printf("Syntax Validated (Found %d jobs in '%s')\n", len(wf.Jobs), wf.Name)

		scanner := security.NewScanner()
		violations := scanner.Scan(wf)

		if len(violations) > 0 {
			fmt.Printf("\nSecurity Policies Enforced: Found %d Violations\n", len(violations))
			for _, v := range violations {
				fmt.Printf("  [%s] %s (Job: %s, Step: %s)\n", v.Severity, v.RuleName, v.JobName, v.StepName)
				fmt.Printf("      -> %s\n", v.Description)
			}
			os.Exit(1)
		} else {
			fmt.Println("Security Policies Enforced (0 Violations)")
		}

		projection := finops.Estimate(wf)

		fmt.Println("\nShift-Left Cost Projections:")
		for _, jc := range projection.Jobs {
			fmt.Printf("  - Job [%s] on '%s' ($%.3f/min)\n", jc.JobName, jc.MachineType, jc.CostPerMinute)
		}
		fmt.Printf("\n  Est. Cost Per Run:   $%.3f\n", projection.TotalEstimated)
		fmt.Printf("  Worst-Case Timeout:  $%.3f\n", projection.TotalWorstCase)
		fmt.Printf("  Est. Monthly Impact: $%.2f (Assuming 100 runs/mo)\n", projection.MonthlyProjected)
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
