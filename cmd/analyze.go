package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Hoaqim/optiflow-cli/pkg/finops"
	"github.com/Hoaqim/optiflow-cli/pkg/security"
	"github.com/Hoaqim/optiflow-cli/pkg/workflow"
	"github.com/spf13/cobra"
)

var customRunners []string

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze a CI/CD workflow file or directory",
	Long: `Parses your YAML workflows, checks for syntax and logic errors,
scans for security vulnerabilities, and projects execution costs.

Example:
  optiflow analyze .github/workflows/main.yml
  optiflow analyze .github/workflows/ --runner "self-hosted=0.01" --runner "gpu-runner=0.75"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := args[0]

		hasSecurityViolations := false
		totalViolationsCount := 0

		for _, cr := range customRunners {
			parts := strings.SplitN(cr, "=", 2)
			if len(parts) == 2 {
				cost, err := strconv.ParseFloat(parts[1], 64)
				if err == nil {
					runnerName := strings.ToLower(parts[0])
					finops.PricingMatrix[runnerName] = cost
					if Verbose {
						fmt.Printf("[DEBUG] Added custom runner: %s at $%.3f/min\n", runnerName, cost)
					}
				} else {
					fmt.Printf("[WARNING] Invalid cost format for runner '%s'. Expected float, got: %s\n", parts[0], parts[1])
				}
			} else {
				fmt.Printf("[WARNING] Invalid --runner format: '%s'. Use 'name=cost' (e.g., 'self-hosted=0.01')\n", cr)
			}
		}

		info, err := os.Stat(targetPath)
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", targetPath, err)
			os.Exit(1)
		}

		var filesToAnalyze []string

		if info.IsDir() {
			if Verbose {
				fmt.Printf("[DEBUG] Scanning directory: %s\n", targetPath)
			}
			err := filepath.WalkDir(targetPath, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && (strings.HasSuffix(d.Name(), ".yml") || strings.HasSuffix(d.Name(), ".yaml")) {
					filesToAnalyze = append(filesToAnalyze, path)
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error walking directory: %v\n", err)
				os.Exit(1)
			}
		} else {
			filesToAnalyze = append(filesToAnalyze, targetPath)
		}

		if len(filesToAnalyze) == 0 {
			fmt.Println("No YAML files found to analyze.")
			os.Exit(0)
		}

		totalViolations := 0
		var globalProjection finops.Projection

		for _, file := range filesToAnalyze {
			if Verbose {
				fmt.Printf("\n[DEBUG] Starting analysis on: %s\n", file)
			}

			wf, err := workflow.Parse(file)
			if err != nil {
				fmt.Printf("\n[ERROR] Syntax Error in %s:\n%v\n", file, err)
				continue
			}

			fmt.Printf("\n--- Analyzing %s ---\n", file)
			fmt.Printf("Syntax Validated (Found %d jobs in '%s')\n", len(wf.Jobs), wf.Name)

			scanner := security.NewScanner()
			violations := scanner.Scan(wf)

			if len(violations) > 0 {
				hasSecurityViolations = true
				totalViolationsCount += len(violations)

				fmt.Printf("\nSecurity Policies Enforced: Found %d Violations\n", len(violations))
				for _, v := range violations {
					fmt.Printf("  [%s] %s\n", v.Severity, v.Description)
				}
			}

			projection := finops.Estimate(wf)
			fmt.Printf("\nShift-Left Cost Projections: $%.3f per run\n", projection.TotalEstimated)

			globalProjection.TotalEstimated += projection.TotalEstimated
			globalProjection.TotalWorstCase += projection.TotalWorstCase
			globalProjection.MonthlyProjected += projection.MonthlyProjected

			fmt.Println("\nShift-Left Cost Projections:")
			for _, jc := range projection.Jobs {
				fmt.Printf("  - Job [%s] on '%s' ($%.3f/min)\n", jc.JobName, jc.MachineType, jc.CostPerMinute)
			}
			fmt.Printf("  Est. Cost Per Run:   $%.3f\n", projection.TotalEstimated)
		}

		if len(filesToAnalyze) > 1 {
			fmt.Printf("\n=== DIRECTORY SCAN SUMMARY ===\n")
			fmt.Printf("Files Analyzed:            %d\n", len(filesToAnalyze))
			fmt.Printf("Total Security Violations: %d\n", totalViolations)
			fmt.Printf("Total Est. Cost Per Run:   $%.3f\n", globalProjection.TotalEstimated)
			fmt.Printf("Total Worst-Case Timeout:  $%.3f\n", globalProjection.TotalWorstCase)
			fmt.Printf("Total Est. Monthly Impact: $%.2f (Assuming 100 runs/mo)\n", globalProjection.MonthlyProjected)
		} else {
			fmt.Printf("\n  Worst-Case Timeout:  $%.3f\n", globalProjection.TotalWorstCase)
			fmt.Printf("  Est. Monthly Impact: $%.2f (Assuming 100 runs/mo)\n", globalProjection.MonthlyProjected)
		}

		if hasSecurityViolations {
			fmt.Printf("\n[FAIL] Scan completed with %d security violations.\n", totalViolationsCount)
			os.Exit(1)
		} else {
			fmt.Println("\n[SUCCESS] Scan completed. No security violations found.")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringSliceVar(&customRunners, "runner", []string{}, "Define custom runner pricing (format: runner-name=cost_per_min, e.g., 'self-hosted=0.015')")
}
