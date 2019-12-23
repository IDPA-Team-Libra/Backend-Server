package logger

import (
	"log"

	"github.com/jbrodriguez/mlog"
)

func SetupLogger(filePath string, rotationSizeInMB int, numberofBackups int) {
	mlog.DefaultFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	mlog.StartEx(mlog.LevelInfo, filePath, rotationSizeInMB*1024*1024, numberofBackups)
	// mlog.Start(mlog.LevelInfo, "")
}

func BackupLogFiles() {
}

type LogLevel int64

const (
	INFO LogLevel = iota
	WARNING
	ERROR
)

func LogMessage(message string, logLevel LogLevel) {
	if logLevel == INFO {
		mlog.Info(message)
	} else if logLevel == WARNING {
		mlog.Warning(message)
	} else {
		mlog.Trace(message)
	}
}
