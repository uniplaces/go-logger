package go_logger

// Logger is logger for logger
type Logger interface {
	ErrorWithFields(err error, fields map[string]interface{})
	WarningWithFields(message string, fields map[string]interface{})
	InfoWithFields(message string, fields map[string]interface{})
	DebugWithFields(message string, fields map[string]interface{})
}
