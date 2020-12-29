package model

// Defaults all options available as defaults in the configuration
type Defaults struct {
	Duration int    `yaml:"duration"`
	Format   string `yaml:"format,omitempty"`
}

// Config all options available in the configuration
type Config struct {
	Author   string   `yaml:"author"`
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
		Author: author,
		Defaults: Defaults{
			Duration: duration,
			Format:   format,
		},
	}
}
