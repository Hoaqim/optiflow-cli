package security

import (
	"strings"

	"github.com/Hoaqim/optiflow-cli/pkg/workflow"
)

type NpmInstallRule struct{}

func (r *NpmInstallRule) Evaluate(jobName string, job workflow.Job) []Violation {
	var violations []Violation

	for _, step := range job.Steps {
		if strings.Contains(step.Run, "npm install") {
			name := step.Name
			if name == "" {
				name = "Unnamed Step"
			}

			violations = append(violations, Violation{
				RuleName:    "Non-Deterministic Install",
				Description: "Found 'npm install'. Use 'npm ci' for secure, deterministic builds.",
				Severity:    SeverityMedium,
				JobName:     jobName,
				StepName:    name,
			})
		}
	}

	return violations
}

type MutableActionRule struct{}

func (r *MutableActionRule) Evaluate(jobName string, job workflow.Job) []Violation {
	var violations []Violation

	for _, step := range job.Steps {
		if step.Uses != "" {
			if strings.HasSuffix(step.Uses, "@main") || strings.HasSuffix(step.Uses, "@master") {
				name := step.Name
				if name == "" {
					name = step.Uses
				}

				violations = append(violations, Violation{
					RuleName:    "Mutable Action Reference",
					Description: "Action uses a mutable branch (" + step.Uses + "). Pin to a specific commit SHA to prevent supply chain attacks.",
					Severity:    SeverityHigh,
					JobName:     jobName,
					StepName:    name,
				})
			}
		}
	}

	return violations
}
