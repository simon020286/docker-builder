package models

// Image struct definition.
type Image struct {
	Path      string   `yaml:"path"`
	Push      bool     `yaml:"push"`
	Platforms []string `yaml:"platforms,flow"`
}
