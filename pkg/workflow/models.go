package workflow

type Workflow struct {
	Name string         `yaml:"name"`
	On   any            `yaml:"on"` // Can be a list, string, or map
	Jobs map[string]Job `yaml:"jobs"`
}

type Job struct {
	RunsOn  string   `yaml:"runs-on"`
	Needs   []string `yaml:"needs,omitempty"` //Can be a single string or list
	Timeout int      `yaml:"timeout-minutes,omitempty"`
	Steps   []Step   `yaml:"steps"`
}

type Step struct {
	Name string            `yaml:"name,omitempty"`
	Uses string            `yaml:"uses,omitempty"`
	Run  string            `yaml:"run,omitempty"`
	Env  map[string]string `yaml:"env,omitempty"`
}
