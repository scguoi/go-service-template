package logc

import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"template/config"

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
	}
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:03:04",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			return frame.Function, path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
		},
	})
}
