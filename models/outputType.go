package models

// OutputType struct definition.
type OutputType string

const (
	// ConsoleOutput redirect output to console.
	ConsoleOutput OutputType = "console"
	// FileOutput redirect output to file.
	FileOutput OutputType = "file"
)
