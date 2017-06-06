package go_logger

// Config defines logger config
type Config struct {
	environment string
	level       string
}

// New returns initialized config
func New(environment, level string) Config {
	return Config{environment: environment, level: level}
}
