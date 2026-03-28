package workflow

import "gopkg.in/yaml.v3"

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalYAML(value *yaml.Node) error {
	var single string
	if err := value.Decode(&single); err == nil {
		*s = []string{single}
		return nil
	}

	var slice []string
	if err := value.Decode(&slice); err == nil {
		*s = slice
		return nil
	}

	return nil
}

type Workflow struct {
	Name        string            `yaml:"name,omitempty"`
	RunName     string            `yaml:"run-name,omitempty"`
	On          any               `yaml:"on"` // Can be string, array, or map
	Env         map[string]string `yaml:"env,omitempty"`
	Permissions map[string]string `yaml:"permissions,omitempty"`
	Concurrency any               `yaml:"concurrency,omitempty"` // Can be string or map
	Jobs        map[string]Job    `yaml:"jobs"`
}

type Job struct {
	Name           string            `yaml:"name,omitempty"`
	Permissions    map[string]string `yaml:"permissions,omitempty"`
	Needs          StringOrSlice     `yaml:"needs,omitempty"` // Can be string or []string
	If             string            `yaml:"if,omitempty"`
	RunsOn         StringOrSlice     `yaml:"runs-on"` // Can be string or []string (e.g., matrix)
	Environment    any               `yaml:"environment,omitempty"`
	Concurrency    any               `yaml:"concurrency,omitempty"`
	Outputs        map[string]string `yaml:"outputs,omitempty"`
	Env            map[string]string `yaml:"env,omitempty"`
	TimeoutMinutes float64           `yaml:"timeout-minutes,omitempty"`
	ContinueError  any               `yaml:"continue-on-error,omitempty"`
	Steps          []Step            `yaml:"steps"`
}

type Step struct {
	Id               string            `yaml:"id,omitempty"`
	If               string            `yaml:"if,omitempty"`
	Name             string            `yaml:"name,omitempty"`
	Uses             string            `yaml:"uses,omitempty"`
	Run              string            `yaml:"run,omitempty"`
	WorkingDirectory string            `yaml:"working-directory,omitempty"`
	Shell            string            `yaml:"shell,omitempty"`
	With             map[string]any    `yaml:"with,omitempty"`
	Env              map[string]string `yaml:"env,omitempty"`
	ContinueError    any               `yaml:"continue-on-error,omitempty"`
	TimeoutMinutes   float64           `yaml:"timeout-minutes,omitempty"`
}
