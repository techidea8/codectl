package logger

type LogLevel = string

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel LogLevel = "panic"
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel LogLevel = "fatal"
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel LogLevel = "error"
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel LogLevel = "warn"
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel LogLevel = "info"
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel LogLevel = "debug"
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel LogLevel = "trace"
)

func TOInt[T int8 | int16 | int32 | int64 | int | uint16 | uint8 | uint32 | uint64](l LogLevel) T {
	var r T = 6
	if l == PanicLevel {
		r = 0
	} else if l == FatalLevel {
		r = 1
	} else if l == ErrorLevel {
		r = 2
	} else if l == WarnLevel {
		r = 3
	} else if l == InfoLevel {
		r = 4
	} else if l == DebugLevel {
		r = 5
	} else if l == TraceLevel {
		r = 6
	} else {
		r = 6
	}
	return r
}
