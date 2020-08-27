package common

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
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

	// Ltimestamp is a flag used to indicate that timestamps should be included in log output
	Ltimestamp = 1 << iota
)

// Logger formats messages to standard out, filtered according to the configured log level.
type Logger struct {
	Level  LogLevel
	Output io.Writer
	flags  int
}

// SetFlags sets the output flags for the logger. Currently, the only flag bit is Ltimestamp.
func (l *Logger) SetFlags(flags int) {
	l.flags = flags
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

	if (l.flags & Ltimestamp) == Ltimestamp {
		fmt.Fprintf(l.Output, "%s ", time.Now().Format(time.RFC3339))
	}

	message := fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...)
	fmt.Fprintf(l.Output, "[%s] %s\n", LevelString(level), message)
}

func (l *Logger) setLevel(level string) {
	switch strings.ToUpper(level) {
	case "ERROR":
		l.Level = ERROR
	case "WARN":
		l.Level = WARN
	case "INFO":
		l.Level = INFO
	case "DEBUG":
		l.Level = DEBUG
	default:
		l.Level = ERROR
		l.Error("Logger assigned unknown log level: %s (defaulting to ERROR)", level)
	}
}

// NewLogger creates new Logger instances.  Valid levels are ERROR, WARN, INFO, and DEBUG.
func NewLogger(level string) *Logger {
	var logger *Logger = &Logger{Output: os.Stdout, flags: 0}

	logger.setLevel(level)

	return logger
}

// NewFileLogger creates new Logger instances.  Valid levels are ERROR, WARN, INFO, and DEBUG.
func NewFileLogger(level, filename string) (*Logger, error) {
	var logger *Logger
	var file *os.File
	var err error

	if file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return nil, fmt.Errorf("unable to open %s for writing: %w", filename, err)
	}

	logger = &Logger{Output: file, flags: Ltimestamp}

	logger.setLevel(level)

	return logger, nil
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
