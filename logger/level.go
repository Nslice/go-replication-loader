package logger

// Level is predefined log levels
type Level string

const (
	// Info indicates that is an inforamtion message
	Info Level = "Info"
	// Error indicates that is an error message
	Error Level = "Error"
	// Fatal indicates that is a message of fatal error
	Fatal Level = "Fatal"
)
