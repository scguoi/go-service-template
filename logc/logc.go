package logc

import (
	"path"
	"runtime"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:03:04",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			return frame.Function, path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
		},
	})
}
