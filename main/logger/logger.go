package logger

func LogLevelToString(level LogLevel) string {
	return string(level)
}

//Logger an interface that will allow us to use multiple ways of outputing messages
type Logger interface {
	Log(string) bool
	GetDepth() int64
	FormatLog(string, LogLevel) string
	CheckState() bool
	GetType() string
}
