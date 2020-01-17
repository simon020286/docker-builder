package models

import "fmt"

// OutputConfig struct definition.
type OutputConfig struct {
	Type OutputType `yaml:"type"`
	Path string     `yaml:"path,omitempty"`
}

// IsValid check if configuration is valid.
func (o *OutputConfig) IsValid() error {
	if o.Type == "" {
		o.Type = ConsoleOutput
	}

	switch o.Type {
	case ConsoleOutput, FileOutput:
	default:
		return fmt.Errorf("Output not valid: %s isn't valid output", o.Type)
	}

	return nil
}
