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

//Logger an interface that will allow us to use multiple ways of outputing messages
type Logger interface {
	Log(string) bool
	GetDepth() int64
	FormatLog(string, LogLevel) string
	CheckState() bool
	GetType() string
}

type ConsoleLogger struct {
}

type FileLogger struct {
}

type RemoteLogger struct {
}
