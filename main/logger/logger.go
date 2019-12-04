package logger

type LogLevel int64

const (
	NORMAL LogLevel = iota
	WARNING
	DANGER
	INFORMATION
	DATABASE
	ROUTINES
	ERRORS
	UNDEFINED
)

func LogLevelToString(level LogLevel) string {
	return string(level)
}

func StringToLogLevel(level string) LogLevel {
	switch level {
	case "normal":
		return NORMAL
	case "warning":
		return WARNING
	case "danger":
		return DANGER
	case "information":
		return INFORMATION
	case "database":
		return DATABASE
	case "routines":
		return ROUTINES
	case "error":
		return ERRORS
	default:
		return UNDEFINED
	}
}

//Logger an interface that will allow us to use multiple ways of outputing messages
type Logger interface {
	Log(string) bool
	GetDepth() int64
	FormatLog(string, LogLevel) string
	CheckState() bool
	GetType() string
}
