package logger

type Level int

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func IsPanicEnabled(logger Logger) bool {
	return logger.IsLevelEnabled(PanicLevel)
}

func IsFatalEnabled(logger Logger) bool {
	return logger.IsLevelEnabled(FatalLevel)
}

func IsErrorEnabled(logger Logger) bool {
	return logger.IsLevelEnabled(ErrorLevel)
}

func IsWarnEnabled(logger Logger) bool {
	return logger.IsLevelEnabled(WarnLevel)
}

func IsInfoEnabled(logger Logger) bool {
	return logger.IsLevelEnabled(InfoLevel)
}

func IsDebugEnabled(logger Logger) bool {
	return logger.IsLevelEnabled(DebugLevel)
}

func IsTraceEnabled(logger Logger) bool {
	return logger.IsLevelEnabled(TraceLevel)
}
