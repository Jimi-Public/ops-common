/*
@Project: ops-common
@Author:  WangChaoQun
@Date:    2023/2/15
@IDE:     GoLand
@File:    gormLog.go
*/

package log

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

var DefaultGormLogs = NewLogrusLogger(DefaultLogs.Log)

// LogrusLogger implements the gorm logger interface using logrus
type LogrusLogger struct {
	Log *logrus.Logger
}

func NewLogrusLogger(log *logrus.Logger) *LogrusLogger {
	return &LogrusLogger{Log: log}
}

// LogMode sets the logger mode
func (logger LogrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := logger
	newLogger.Log.Level = levelToLogrusLevel(level)
	return newLogger
}

// Info prints information about the statement
func (logger LogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.Log.WithContext(ctx).WithFields(logrus.Fields{
		"data": data,
	}).Info(msg)
}

// Warn prints warning about the statement
func (logger LogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.Log.WithContext(ctx).WithFields(logrus.Fields{
		"data": data,
	}).Warn(msg)
}

// Error prints error about the statement
func (logger LogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.Log.WithContext(ctx).WithFields(logrus.Fields{
		"data": data,
	}).Error(msg)
}

// Trace prints sql statement if log level is greater than debug
func (logger LogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if logger.Log.Level >= logrus.DebugLevel {
		elapsed := time.Since(begin)
		sql, rows := fc()
		logger.Log.WithContext(ctx).WithFields(logrus.Fields{
			"sql":    sql,
			"rows":   rows,
			"time":   fmt.Sprintf("%v", elapsed),
			"caller": getCallerFunctionName(),
			"error":  err,
		}).Debug("SQL")
	}
}

// Convert Gorm LogLevel to Logrus LogLevel
func levelToLogrusLevel(level logger.LogLevel) logrus.Level {
	switch level {
	case logger.Silent:
		return logrus.PanicLevel
	case logger.Error:
		return logrus.ErrorLevel
	case logger.Warn:
		return logrus.WarnLevel
	case logger.Info:
		return logrus.InfoLevel
	default:
		return logrus.DebugLevel
	}
}

// Get the name of the calling function
func getCallerFunctionName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
