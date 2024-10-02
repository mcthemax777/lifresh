package lflog

import (
	"github.com/fluent/fluent-logger-golang/fluent"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"lifresh/define"
	"time"
)

const (
	LogLevelDebug = 1
	LogLevelInfo  = 2
	LogLevelWarn  = 3
	LogLevelError = 4
	LogLevelFatal = 5
	LogLevelPanic = 6
)

var logger *zap.Logger
var fluentLogger *fluent.Fluent
var logError error
var isActiveFluentd bool

func init() {

	logFile := ""
	// initialize the rotator
	if define.OsType == define.OsTypeWindows {
		logFile = "./log/app-%Y-%m-%d-%H.log"
		isActiveFluentd = false
	} else {
		logFile = "/var/log/app-%Y-%m-%d-%H.log"
		isActiveFluentd = false

	}

	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(60*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		panic(err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""

	// add the encoder config and rotator to create a new zap logger
	w := zapcore.AddSync(rotator)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		zap.InfoLevel)
	logger = zap.New(core)

	if isActiveFluentd {
		////fluentd
		fluentLogger, logError = fluent.New(fluent.Config{FluentNetwork: "tcp", FluentHost: "host.docker.internal"})
		if logError != nil {
			panic(logError.Error())
		}
	}
}

func InfoLogging(msg string) {
	Logging(LogLevelInfo, msg)
}

func Logging(level int, msg string) {

	switch level {
	case LogLevelDebug:
		logger.Debug(msg)
	case LogLevelInfo:
		logger.Info(msg)
	case LogLevelWarn:
		logger.Warn(msg)
	case LogLevelError:
		logger.Error(msg)
	case LogLevelFatal:
		logger.Fatal(msg)
	case LogLevelPanic:
		logger.Panic(msg)
	default:
		logger.Panic(msg)
	}

	if isActiveFluentd {
		mapStringData := map[string]string{
			"msg": msg,
		}

		switch level {
		case LogLevelDebug:
			fluentLogger.Post("debug", mapStringData)
		case LogLevelInfo:
			fluentLogger.Post("info", mapStringData)
		case LogLevelWarn:
			fluentLogger.Post("warn", mapStringData)
		case LogLevelError:
			fluentLogger.Post("error", mapStringData)
		case LogLevelFatal:
			fluentLogger.Post("fatal", mapStringData)
		case LogLevelPanic:
			fluentLogger.Post("panic", mapStringData)
		default:
			fluentLogger.Post("panic", mapStringData)
		}
	}
}
