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
)

func LogLevelToString(level LogLevel) string {
	return ""
}

//Logger an interface that will allow us to use multiple ways of outputing messages
type Logger interface {
	Log(string) bool
	GetDepth() int64
	FormatLog(string, LogLevel) string
	CheckState() bool
	GetType() string
}
