package logger

// Fields represents a map of fields
type Fields map[string]interface{}

type extraField struct {
	key   string
	value interface{}
}

// Logger is logger for logger
type Logger interface {
	ErrorWithFields(message string, fields Fields)
	Error(message string)

	WarningWithFields(message string, fields Fields)
	Warning(message string)

	InfoWithFields(message string, fields Fields)
	Info(message string)

	DebugWithFields(message string, fields Fields)
	Debug(message string)
}
