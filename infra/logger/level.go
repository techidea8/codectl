package logger

type LogLevel = uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel LogLevel = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = PanicLevel + 1
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = PanicLevel + 2
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = PanicLevel + 3
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = PanicLevel + 4
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = PanicLevel + 4
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = PanicLevel + 5
)
