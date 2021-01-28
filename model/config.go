package model

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Defaults all options available as defaults in the configuration
type Defaults struct {
	Author   string `yaml:"author"`
	Duration int    `yaml:"duration"`
	Format   string `yaml:"format,omitempty"`
}

// Config all options available in the configuration
type Config struct {
	Defaults Defaults `yaml:"default"`
}

// NewConfig is the generator for configuration
func NewConfig(author, format string, duration int) *Config {
	if format != "pretty" &&
		format != "json" &&
		format != "yaml" &&
		format != "yml" {
		format = ""
	}
	return &Config{
		Defaults: Defaults{
			Author:   author,
			Duration: duration,
			Format:   format,
		},
	}
}

// WriteYAML takes a writer and outputs a YAML representation of Config to it
func (c *Config) WriteYAML(writer io.Writer) error {
	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	_, err = writer.Write(b)
	return err
}
