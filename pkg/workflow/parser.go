package workflow

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Parse(filePath string) (*Workflow, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var wf Workflow

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	decoder.KnownFields(true)

	if err := decoder.Decode(&wf); err != nil {
		return nil, fmt.Errorf("strict schema validation failed: %w", err)
	}

	return &wf, nil
}
