package models

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

var (
	file *os.File
)

// Config struct definition.
type Config struct {
	Path   string           `yaml:"path"`
	Images map[string]Image `yaml:"images,omitempty"`
	Output OutputConfig     `yaml:"output"`
}

// IsValid check if configuration is valid.
func (s *Config) IsValid() error {
	err := s.Output.IsValid()
	if err != nil {
		return err
	}

	return nil
}

// Close and release memory.
func (s *Config) Close() {
	log.Println("Closing")
	if file != nil {
		err := file.Close()
		if err != nil {
			log.Panicf("Close log file: %s", err)
		}
	}
}

// InitConfig from file.
func InitConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "./config.yaml"
	}

	reader, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	config := &Config{}
	decoder := yaml.NewDecoder(reader)

	err = decoder.Decode(config)

	if err != nil {
		return nil, err
	}

	if err := config.IsValid(); err != nil {
		return nil, err
	}

	if err := initLog(config.Output); err != nil {
		return nil, err
	}

	return config, nil
}

func initLog(config OutputConfig) error {
	if config.Type == FileOutput {
		file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		log.SetOutput(file)
	}

	return nil
}
