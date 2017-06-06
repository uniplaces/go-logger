package go_logger

// Config defines logger config
type Config struct {
	environment string
	level       string
}

// NewConfig returns initialized config
func NewConfig(environment, level string) Config {
	return Config{environment: environment, level: level}
}
