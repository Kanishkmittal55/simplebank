package worker

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.print(zerolog.DebugLevel, args...)
}

// Info logs a message at Info level.
func (logger *Logger) Info(args ...interface{}) {
	logger.print(zerolog.InfoLevel, args...)
}

// Warn logs a message at Warning level.
func (logger *Logger) Warn(args ...interface{}) {
	logger.print(zerolog.WarnLevel, args...)
}

// Error logs a message at Error level.
func (logger *Logger) Error(args ...interface{}) {
	logger.print(zerolog.ErrorLevel, args...)
}

// Fatal logs a message at Fatal level
// and process will exit with status set to 1.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.print(zerolog.FatalLevel, args...)
}
