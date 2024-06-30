package logger

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

type LogrusImpl struct {
	logger *logrus.Logger
}

func NewLogrus(filepath string, level logrus.Level) LogrusImpl {

	// 创建一个新的 logger
	loggerus := logrus.New()
	w1 := os.Stdout
	// 设置日志级别
	loggerus.SetLevel(level)
	// 创建一个文件句柄，用于写入日志
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 设置日志输出到文件
	loggerus.SetOutput(io.MultiWriter(w1, file))
	// 配置日志格式
	loggerus.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	return LogrusImpl{
		logger: loggerus,
	}
}
func (l LogrusImpl) SetLevel(level LogLevel) {
	l.logger.SetLevel(logrus.Level(level))
}
func (l LogrusImpl) Debug(arg ...any) {
	l.logger.Debug(arg...)
}
func (l LogrusImpl) Info(args ...any) {
	l.logger.Info(args...)
}

func (l LogrusImpl) Error(args ...any) {
	l.logger.Error(args...)
}

func (l LogrusImpl) Debugf(fmt string, args ...any) {
	l.logger.Debugf(fmt, args...)
}

func (l LogrusImpl) Infof(fmt string, args ...any) {
	l.logger.Infof(fmt, args...)
}
func (l LogrusImpl) Errorf(fmt string, args ...any) {
	l.logger.Errorf(fmt, args...)
}
