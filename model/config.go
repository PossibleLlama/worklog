package model

// Defaults all options available as defaults in the configuration
type Defaults struct {
	Duration int
}

// Config all options available in the configuration
type Config struct {
	Author   string   `yaml:"author"`
	Defaults Defaults `yaml:"default"`
}

// NewConfig is the generator for configuration
func NewConfig(author string, duration int) *Config {
	return &Config{
		Author: author,
		Defaults: Defaults{
			Duration: duration,
		},
	}
}
