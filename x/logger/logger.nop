package logger

import (
	"io"

	"github.com/sirupsen/logrus"
	"github.com/wenewzhang/core/logger"
)

type LoggerOptions struct {
}

type LoggerOption func(opts *LoggerOptions)

func OutputLoggerOption(out io.Writer) LoggerOption {
	return func(opts *LoggerOptions) {
	}
}

func FormatLoggerOption(format logger.LogFormat) LoggerOption {
	return func(opts *LoggerOptions) {
	}
}

func LevelLoggerOption(level logger.LogLevel) LoggerOption {
	return func(opts *LoggerOptions) {
	}
}

type logrusLogger struct {
}

func NewLogger(opts ...LoggerOption) logger.Logger {
	return &logrusLogger{}
}

// WithFields adds new fields to log.
func (l *logrusLogger) WithFields(fields map[string]any) logger.Logger {
	return &logrusLogger{}
}

// Trace logs a message at level Trace.
func (l *logrusLogger) Trace(args ...any) {
}

// Tracef logs a message at level Trace.
func (l *logrusLogger) Tracef(format string, args ...any) {
}

// Debug logs a message at level Debug.
func (l *logrusLogger) Debug(args ...any) {
}

// Debugf logs a message at level Debug.
func (l *logrusLogger) Debugf(format string, args ...any) {
}

// Info logs a message at level Info.
func (l *logrusLogger) Info(args ...any) {
}

// Infof logs a message at level Info.
func (l *logrusLogger) Infof(format string, args ...any) {
}

// Warn logs a message at level Warn.
func (l *logrusLogger) Warn(args ...any) {
}

// Warnf logs a message at level Warn.
func (l *logrusLogger) Warnf(format string, args ...any) {
}

// Error logs a message at level Error.
func (l *logrusLogger) Error(args ...any) {
}

// Errorf logs a message at level Error.
func (l *logrusLogger) Errorf(format string, args ...any) {
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
func (l *logrusLogger) Fatal(args ...any) {
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func (l *logrusLogger) Fatalf(format string, args ...any) {
}

func (l *logrusLogger) GetLevel() logger.LogLevel {
	return ""
}

func (l *logrusLogger) IsLevelEnabled(level logger.LogLevel) bool {
	return false
}

func (l *logrusLogger) log(level logrus.Level, args ...any) {
}

func (l *logrusLogger) logf(level logrus.Level, format string, args ...any) {
}

func (l *logrusLogger) caller(skip int) string {
	return ""
}
