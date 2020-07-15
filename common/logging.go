package common

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// LogLevel is base type of log levels
type LogLevel int

const (
	// ERROR is a log level
	ERROR LogLevel = 3
	// WARN is a log level
	WARN LogLevel = 2
	// INFO is a log level
	INFO LogLevel = 1
	// DEBUG is a log level
	DEBUG LogLevel = 0
)

// Logger formats messages to standard out, filtered according to the configured log level.
type Logger struct {
	Level  LogLevel
	Output io.Writer
}

// Error logs messages at the ERROR level
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// Warn logs messages at the WARN level
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

// Info logs messages at the INFO level
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Debug logs messages at the DEBUG level
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level < l.Level {
		return
	}

	message := fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...)
	fmt.Fprintf(l.Output, "[%s] %s\n", LevelString(level), message)
}

// NewLogger creates new Logger instances.  Valid levels are ERROR, WARN, INFO, and DEBUG.
func NewLogger(level string) *Logger {
	var logger *Logger = &Logger{Output: os.Stdout}

	switch strings.ToUpper(level) {
	case "ERROR":
		logger.Level = ERROR
	case "WARN":
		logger.Level = WARN
	case "INFO":
		logger.Level = INFO
	case "DEBUG":
		logger.Level = DEBUG
	default:
		logger.Level = ERROR
		logger.Error("Logger initialized with unknown log level: %s (defaulting to ERROR)", level)
	}

	return logger
}

// LevelString maps a LogLevel type to its string value
func LevelString(level LogLevel) string {
	switch level {
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
