package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents the severity level of log messages
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger represents a colorful terminal logger
type Logger struct {
	mu         sync.Mutex
	output     io.Writer
	level      LogLevel
	showCaller bool
	colors     struct {
		debug *color.Color
		info  *color.Color
		warn  *color.Color
		error *color.Color
		fatal *color.Color
		timestamp *color.Color
		caller    *color.Color
		message   *color.Color
	}
}

// New creates a new logger instance
func New() *Logger {
	logger := &Logger{
		output:     os.Stdout,
		level:      INFO,
		showCaller: true,
	}

	// Initialize colors
	logger.colors.debug = color.New(color.FgHiCyan, color.Bold)
	logger.colors.info = color.New(color.FgHiGreen, color.Bold)
	logger.colors.warn = color.New(color.FgHiYellow, color.Bold)
	logger.colors.error = color.New(color.FgHiRed, color.Bold)
	logger.colors.fatal = color.New(color.FgHiMagenta, color.Bold, color.BgBlack)
	logger.colors.timestamp = color.New(color.FgHiWhite, color.Faint)
	logger.colors.caller = color.New(color.FgHiBlue, color.Italic)
	logger.colors.message = color.New(color.FgHiWhite)

	return logger
}

// SetLevel sets the minimum log level to output
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the output writer
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetShowCaller enables or disables caller information
func (l *Logger) SetShowCaller(show bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showCaller = show
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
	os.Exit(1)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Get timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Get caller information
	var callerInfo string
	if l.showCaller {
		_, file, line, ok := runtime.Caller(2) // 2 levels up the call stack
		if ok {
			// Shorten file path for better readability
			parts := strings.Split(file, "/")
			if len(parts) > 2 {
				callerInfo = fmt.Sprintf("%s/%s:%d", 
					parts[len(parts)-2], parts[len(parts)-1], line)
			} else {
				callerInfo = fmt.Sprintf("%s:%d", file, line)
			}
		}
	}

	// Format the message
	message := fmt.Sprintf(format, v...)

	// Get level string and color
	var levelStr string
	var levelColor *color.Color

	switch level {
	case DEBUG:
		levelStr = "DEBUG"
		levelColor = l.colors.debug
	case INFO:
		levelStr = "INFO"
		levelColor = l.colors.info
	case WARN:
		levelStr = "WARN"
		levelColor = l.colors.warn
	case ERROR:
		levelStr = "ERROR"
		levelColor = l.colors.error
	case FATAL:
		levelStr = "FATAL"
		levelColor = l.colors.fatal
	}

	// Build the log line
	var logLine strings.Builder

	// Timestamp
	l.colors.timestamp.Fprintf(&logLine, "[%s] ", timestamp)

	// Level
	levelColor.Fprintf(&logLine, "%-5s ", levelStr)

	// Caller info (if enabled)
	if l.showCaller && callerInfo != "" {
		l.colors.caller.Fprintf(&logLine, "(%s) ", callerInfo)
	}

	// Message
	l.colors.message.Fprintf(&logLine, "%s", message)

	// Write to output
	fmt.Fprintln(l.output, logLine.String())
}

// Simple global logger instance for convenience
var defaultLogger = New()

// Package-level functions for easy access
func Debug(format string, v ...interface{}) {
	defaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	defaultLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	defaultLogger.Fatal(format, v...)
}

func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

func SetShowCaller(show bool) {
	defaultLogger.SetShowCaller(show)
}