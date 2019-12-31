package logger

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger(filePath string, rotationSizeInMB int, numberofBackups int) {
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

func SyncLogger() {
	loggerInstance.Sync()
}

var loggerInstance *zap.SugaredLogger

type LogLevel int64

const (
	INFO LogLevel = iota
	WARNING
	ERROR
)

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
func LogWithoutFields(message string, logLevel LogLevel) {
	if logLevel == INFO {
		loggerInstance.Info(message)
	} else if logLevel == WARNING {
		loggerInstance.Warn(message)
	} else {
		loggerInstance.Error(message)
	}
}
func LogMessageWithOrigin(message string, logLevel LogLevel, origin string) {
	message = fmt.Sprintf("%s -- %s", message, origin)
	LogMessage(message, logLevel)
}
