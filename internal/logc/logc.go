package logc

import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"template/internal/config"
	"time"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func init() {
	switch strings.ToLower(config.ServiceConfig.LogLevel) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:05:04",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			return frame.Function, path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
		},
	})
	if strings.ToLower(config.ServiceConfig.LogOut) == "file" {
		writer, _ := rotateLogs.New(
			config.ServiceConfig.LogFile+".%Y%m%d",
			rotateLogs.WithLinkName(config.ServiceConfig.LogFile),
			rotateLogs.WithMaxAge(time.Duration(config.ServiceConfig.LogMaxDays*24)*time.Hour),
			rotateLogs.WithRotationTime(time.Duration(1*24)*time.Hour),
		)
		log.SetOutput(writer)
	}
}
