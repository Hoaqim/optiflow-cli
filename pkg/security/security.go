package security

import (
	"github.com/Hoaqim/optiflow-cli/pkg/workflow"
)

type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
	SeverityLow    Severity = "LOW"
)

type Violation struct {
	RuleName    string
	Description string
	Severity    Severity
	JobName     string
	StepName    string
}

type Rule interface {
	Evaluate(jobName string, job workflow.Job) []Violation
}

type Scanner struct {
	rules []Rule
}

func NewScanner() *Scanner {
	return &Scanner{
		rules: []Rule{
			&NpmInstallRule{},
			&MutableActionRule{},
		},
	}
}

func (s *Scanner) Scan(wf *workflow.Workflow) []Violation {
	var allViolations []Violation

	for jobName, job := range wf.Jobs {
		for _, rule := range s.rules {
			violations := rule.Evaluate(jobName, job)
			allViolations = append(allViolations, violations...)
		}
	}

	return allViolations
}
