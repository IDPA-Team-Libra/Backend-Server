package logger

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//SetupLogger creates a logger and sets the ouput file for the log-location
func SetupLogger(filePath string) {
	layout := "2006-01-02"
	date := time.Now().Local()
	t := date.Format(layout)
	filePath = fmt.Sprintf("%s-%s.log", filePath, t)
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	snycWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    50, // megabytes
		MaxBackups: 10,
		MaxAge:     60,    //days
		Compress:   false, // disabled by default
	})
	loggerInstance = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		snycWriter,
		atom,
	)).Sugar()
}

//SyncLogger closes logger and flushes all hanging entries
func SyncLogger() {
	loggerInstance.Sync()
}

var loggerInstance *zap.SugaredLogger

//LogLevel iota for different log levels
type LogLevel int64

const (
	//INFO -- default log level
	INFO LogLevel = iota
	//WARNING -- log level for events that should be investigated
	WARNING
	//ERROR -- something happend, that should be solved asap
	ERROR
)

//LogMessage logs a message to the file
func LogMessage(message string, logLevel LogLevel, fields ...zap.Field) {
	if len(fields) == 0 {
		LogWithoutFields(message, logLevel)
		return
	}
	if logLevel == INFO {
		loggerInstance.Info(message, fields)
	} else if logLevel == WARNING {
		loggerInstance.Warn(message, fields)
	} else {
		loggerInstance.Error(message, fields)
	}
}

//LogWithoutFields logs a message without zap fields
func LogWithoutFields(message string, logLevel LogLevel) {
	if logLevel == INFO {
		loggerInstance.Info(message)
	} else if logLevel == WARNING {
		loggerInstance.Warn(message)
	} else {
		loggerInstance.Error(message)
	}
}

//LogMessageWithOrigin logs a message with its origin
func LogMessageWithOrigin(message string, logLevel LogLevel, origin string) {
	message = fmt.Sprintf("%s -- %s", message, origin)
	LogMessage(message, logLevel)
}
