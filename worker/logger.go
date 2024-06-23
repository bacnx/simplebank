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

func print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Debug(args ...interface{}) {
	print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	print(zerolog.InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	print(zerolog.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	print(zerolog.FatalLevel, args...)
}
