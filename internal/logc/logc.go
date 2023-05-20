package logc

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"path"
	"runtime"
	"strconv"
	"strings"
	"template/internal/config"
	"time"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

var producer *kafka.Producer // kafka producer
var chain chan kafka.Event   // 异步发送缓存

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
		TimestampFormat: "2006-01-02 15:04:05",
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

	if config.ServiceConfig.IsRemoteBizLog {
		conf := buildConfigMap(config.ServiceConfig.KafkaClientExt)
		p, err := kafka.NewProducer(&conf)
		if err != nil {
			log.Errorf("kafka new producer failed, err: %v", err)
		}
		producer = p
		chain = make(chan kafka.Event, config.ServiceConfig.AsyncChainSize)
	}
}

func buildConfigMap(configMap map[string]interface{}) kafka.ConfigMap {
	conf := kafka.ConfigMap{}
	for k, v := range configMap {
		if c, ok := v.(map[string]interface{}); ok {
			conf[k] = buildConfigMap(c)
		} else {
			conf[k] = v
		}
	}
	return conf
}
