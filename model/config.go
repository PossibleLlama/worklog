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

type Repo struct {
	Type string `yaml:"type"`
	Path string `yaml:"path,omitempty"`
}

// Config all options available in the configuration
type Config struct {
	Defaults Defaults `yaml:"default"`
	Repo     Repo     `yaml:"repo"`
}

// NewConfig is the generator for configuration
func NewConfig(def Defaults, repo Repo) *Config {
	if def.Format != "pretty" &&
		def.Format != "json" &&
		def.Format != "yaml" &&
		def.Format != "yml" {
		def.Format = ""
	}
	if repo.Type != "bolt" &&
		repo.Type != "legacy" {
		repo.Type = ""
	}
	return &Config{
		Defaults: def,
		Repo:     repo,
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
