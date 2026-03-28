package finops

import (
	"fmt"
	"strings"

	"github.com/Hoaqim/optiflow-cli/pkg/workflow"
)

// Prices in USD
var PricingMatrix = map[string]float64{
	"ubuntu-latest":  0.008,
	"ubuntu-22.04":   0.008,
	"ubuntu-20.04":   0.008,
	"windows-latest": 0.016,
	"windows-2022":   0.016,
	"macos-latest":   0.080,
	"macos-13":       0.080,
	"macos-14":       0.080,
}

const (
	DefaultGitHubTimeout = 360.0
	DefaultAverageRun    = 5.0
)

type JobCost struct {
	JobName       string
	MachineType   string
	CostPerMinute float64
	EstimatedCost float64
	WorstCaseCost float64
}

type Projection struct {
	Jobs             []JobCost
	TotalEstimated   float64
	TotalWorstCase   float64
	MonthlyProjected float64 // Assuming 100 runs a month
}

func Estimate(wf *workflow.Workflow) Projection {
	proj := Projection{}

	for jobName, job := range wf.Jobs {
		var machineStr string

		if len(job.RunsOn) > 0 {
			machineStr = job.RunsOn[0]
		} else {
			machineStr = "ubuntu-latest"
		}
		machine := strings.ToLower(machineStr)

		costPerMin, exists := PricingMatrix[machine]
		if !exists {
			fmt.Printf("Runner unknown! Set it up with --runner. Defaulting to ubuntu")
			costPerMin = 0.008
		}

		timeout := job.TimeoutMinutes
		if timeout == 0 {
			timeout = DefaultGitHubTimeout
		}

		estCost := costPerMin * DefaultAverageRun
		worstCost := costPerMin * timeout

		proj.Jobs = append(proj.Jobs, JobCost{
			JobName:       jobName,
			MachineType:   machineStr,
			CostPerMinute: costPerMin,
			EstimatedCost: estCost,
			WorstCaseCost: worstCost,
		})

		proj.TotalEstimated += estCost
		proj.TotalWorstCase += worstCost
	}

	proj.MonthlyProjected = proj.TotalEstimated * 100

	return proj
}
