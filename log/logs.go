/*
@Project: ops-common
@Author:  WangChaoQun
@Date:    2023/2/10
@IDE:     GoLand
@File:    logs.go
*/

package log

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const TraceName = "trace_id"

var DefaultLogs = NewLog()

type Log struct {
	Log *logrus.Logger
}

//func (l *Log) WithContext(ctx *gin.Context) *Log {
//	l.Log.WithContext(ctx)
//	return l
//}

func NewLog() *Log {
	var jsonFormatter = logrus.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000",
		DisableTimestamp: false,
		//DisableHTMLEscape: false,
		//DataKey:           "",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
		CallerPrettyfier: nil,
		//PrettyPrint:      false,
	}
	logPath, _ := os.Getwd() // 后期可以配置
	logName := fmt.Sprintf("%s/logs/access.Log", logPath)
	writer, _ := rotatelogs.New(logName+"_%Y%m%d",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),
		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour*24),
		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		// rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(3),
	)
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.DebugLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.PanicLevel: writer,
		logrus.FatalLevel: writer,
	}, &jsonFormatter)
	Logs := logrus.New()
	Logs.SetFormatter(&jsonFormatter)
	Logs.SetOutput(os.Stdout)
	Logs.AddHook(lfHook)
	Logs.SetReportCaller(true)
	return &Log{Log: Logs}
}

func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(TraceName)
		if traceID == "" {
			traceID = uuid.NewV4().String()
		}
		c.Set(TraceName, traceID)
		DefaultLogs.Log.AddHook(NewLoggerHook(traceID))
		c.Next()
	}
}

// TraceLoggerHook Hook log 上下取出trace_id
type TraceLoggerHook struct {
	Trace string
}

func (h *TraceLoggerHook) Fire(entry *logrus.Entry) error {
	entry.Data[TraceName] = h.Trace
	return nil
}

func (h *TraceLoggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func NewLoggerHook(t string) logrus.Hook {
	hook := &TraceLoggerHook{
		Trace: t,
	}
	return hook
}
