package go_logger

// Logger is logger for logger
type Logger interface {
	ErrorWithFields(message string, fields map[string]interface{})
	Error(message string)

	WarningWithFields(message string, fields map[string]interface{})
	Warning(message string)

	InfoWithFields(message string, fields map[string]interface{})
	Info(message string)

	DebugWithFields(message string, fields map[string]interface{})
	Debug(message string)
}
