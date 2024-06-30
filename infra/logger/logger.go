package logger

type ILogger interface {
	Debug(arg ...any)
	Debugf(fmt string, arg ...any)
	Error(arg ...any)
	Errorf(fmt string, arg ...any)
	Info(arg ...any)
	Infof(fmt string, arg ...any)
	SetLevel(LogLevel)
}
