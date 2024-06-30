package logger

var DefaultLogger ILogger = NewLogrus("app.log", InfoLevel)

func Use(logger ILogger) {
	DefaultLogger = logger
}

func SetLevel(level LogLevel) {
	DefaultLogger.SetLevel(level)
}
func Debug(arg ...any) {
	DefaultLogger.Debug(arg...)
}
func Info(arg ...any) {
	DefaultLogger.Info(arg...)
}

func Error(arg ...any) {
	DefaultLogger.Error(arg...)
}

func Debugf(fmt string, arg ...any) {
	DefaultLogger.Debugf(fmt, arg...)
}

func Infof(fmt string, arg ...any) {
	DefaultLogger.Infof(fmt, arg...)
}
func Errorf(fmt string, arg ...any) {
	DefaultLogger.Errorf(fmt, arg...)
}
