package logger

import (
	"io"
	"log"
	"os"
)

type std struct {
	logger *log.Logger
	level  LogLevel
}

func Std() *std {
	return NewStd(InfoLevel)
}
func NewStd(level LogLevel, writers ...io.Writer) *std {
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}
	writer := io.MultiWriter(writers...)
	logger := log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &std{
		level:  level,
		logger: logger,
	}
}

func (l std) Debug(args ...any) {
	if TOInt[int32](l.level) > TOInt[int32](DebugLevel) {
		return
	}
	l.logger.Print(args...)
}
func (l std) Debugf(fmt string, args ...any) {
	if TOInt[int32](l.level) > TOInt[int32](DebugLevel) {
		return
	}
	l.logger.Printf(fmt, args...)
}

func (l *std) SetLevel(level LogLevel) {
	l.level = level
}
func (l std) Info(args ...any) {
	if TOInt[int32](l.level) > TOInt[int32](InfoLevel) {
		return
	}
	l.logger.Print(args...)
}
func (l std) Infof(fmt string, args ...any) {
	if TOInt[int32](l.level) > TOInt[int32](InfoLevel) {
		return
	}
	l.logger.Printf(fmt, args...)
}

func (l std) Error(args ...any) {
	if TOInt[int32](l.level) > TOInt[int32](ErrorLevel) {
		return
	}
	l.logger.Print(args...)
}

func (l std) Errorf(fmt string, args ...any) {
	if TOInt[int32](l.level) > TOInt[int32](ErrorLevel) {
		return
	}
	l.logger.Printf(fmt, args...)
}
