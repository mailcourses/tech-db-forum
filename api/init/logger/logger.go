package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type LogrusLogger struct {
	Logrus *logrus.Entry
}

var GlobalLogger *LogrusLogger

func (logger *LogrusLogger) JsonLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		logger.Logrus.Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   "",
			DisableTimestamp:  false,
			DisableHTMLEscape: false,
			DataKey:           "now",
			FieldMap:          nil,
			CallerPrettyfier:  nil,
			PrettyPrint:       true,
		})
		logger.Logrus.WithFields(logrus.Fields{
			"method":      ctx.Request().Method,
			"remote_addr": ctx.Request().RemoteAddr,
			"work_time":   time.Since(start),
		}).Debug(ctx.Request().URL.Path)
		return next(ctx)
	}
}

func (logger *LogrusLogger) ColoredLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		logger.Logrus.Logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:               false,
			DisableColors:             false,
			ForceQuote:                false,
			DisableQuote:              false,
			EnvironmentOverrideColors: false,
			DisableTimestamp:          false,
			FullTimestamp:             false,
			TimestampFormat:           "",
			DisableSorting:            false,
			SortingFunc:               nil,
			DisableLevelTruncation:    true,
			PadLevelText:              true,
			QuoteEmptyFields:          false,
			FieldMap:                  nil,
			CallerPrettyfier:          nil,
		})
		logger.Logrus.WithFields(logrus.Fields{
			"method":      ctx.Request().Method,
			"remote_addr": ctx.Request().RemoteAddr,
			"work_time":   time.Since(start),
		}).Warn(ctx.Request().URL.Path)
		return next(ctx)
	}
}

func InitLogrus(port string, dbType string) (*LogrusLogger, error) {
	host, _ := os.Hostname()
	contextLogger := logrus.WithFields(logrus.Fields{
		"logger":   "LOGRUS",
		"host":     host,
		"port":     port,
		"database": dbType,
	})
	contextLogger.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    true,
		PadLevelText:              true,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})
	contextLogger.Logger.SetLevel(logrus.TraceLevel)
	GlobalLogger = new(LogrusLogger)
	GlobalLogger.Logrus = contextLogger
	return GlobalLogger, nil
}
