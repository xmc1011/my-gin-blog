package logger

import (
	"github.com/sirupsen/logrus"
	c "my-blog/internal/config"
	"os"
	"time"
)

var logger = logrus.New()

// Define logrus alias
var (
	Tracef          = logrus.Tracef
	Debugf          = logrus.Debugf
	Infof           = logrus.Infof
	Warnf           = logrus.Warnf
	Errorf          = logrus.Errorf
	Fatalf          = logrus.Fatalf
	Panicf          = logrus.Panicf
	Printf          = logrus.Printf
	SetOutput       = logrus.SetOutput
	SetReportCaller = logrus.SetReportCaller
	StandardLogger  = logrus.StandardLogger
	ParseLevel      = logrus.ParseLevel
)

type Level = logrus.Level

// Define logger level
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// SetLevel Set logger level
func SetLevel(level Level) {
	logrus.SetLevel(level)
}

func InitLogger(conf *c.Config) *logrus.Logger {

	// 设置日志级别
	switch conf.Log.Level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	//设置日志格式
	switch conf.Log.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	case "text":
		fallthrough
	default:
		// 使用文本格式
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})

	}
	// 将日志输出到标准输出
	logger.SetOutput(os.Stdout)

	// 设置为全局日志记录器
	logrus.SetFormatter(logger.Formatter)
	logrus.SetOutput(logger.Out)
	logrus.SetLevel(logger.GetLevel())

	return logger
}
