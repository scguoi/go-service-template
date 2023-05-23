package logc

import (
	"fmt"
	"strings"
	"template/internal/config"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type BizLog struct {
	logType    string            // 日志类型
	cpName     string            // 组件名字
	stm        time.Time         // 开始时间
	etm        time.Time         // 结束时间
	utm        int32             // 耗时 ms
	raw        []string          // 原始日志
	structured map[string]string // kv结构化日志
}

func NewBizLog() *BizLog {
	b := &BizLog{
		logType:    config.ServiceConfig.LogType,
		cpName:     config.ServiceConfig.CpName,
		stm:        time.Now(),
		raw:        make([]string, 0),
		structured: make(map[string]string),
	}
	b.structured["logtype"] = b.logType
	b.structured["cpname"] = b.cpName
	b.structured["stm"] = b.stm.Format("2006-01-02 15:04:05.000")
	return b
}

func (b *BizLog) LoggerRaw(raw string) {
	b.raw = append(b.raw, raw)
}

func (b *BizLog) Printf(format string, args ...interface{}) {
	b.raw = append(b.raw, fmt.Sprintf(format, args...))
}

func (b *BizLog) Print(args ...interface{}) {
	str := fmt.Sprintln(args...)
	str = strings.TrimSuffix(str, "\n")
	b.raw = append(b.raw, str)
}

func (b *BizLog) LoggerStructured(key string, value string) {
	b.structured[key] = value
}

func (b *BizLog) LoggerStructuredBatch(logs map[string]string) {
	for k, v := range logs {
		b.structured[k] = v
	}
}

func (b *BizLog) LoggerEnd() {
	b.structured["inner"] = "[" + strings.Join(b.raw, " ") + "]"
	b.etm = time.Now()
	b.utm = int32(b.etm.Sub(b.stm).Milliseconds())
	b.structured["etm"] = b.etm.Format("2006-01-02 15:04:05.000")
	b.structured["utm"] = fmt.Sprintf("%d", b.utm)
	if _, ok := b.structured["traceid"]; !ok {
		b.structured["traceid"] = uuid.New().String()
	}
	if _, ok := b.structured["appid"]; !ok {
		b.structured["appid"] = config.ServiceConfig.DefaultAppID
	}
	if _, ok := b.structured["productline"]; !ok {
		b.structured["productline"] = config.ServiceConfig.DefaultProductLine
	}
	if _, ok := b.structured["bizid"]; !ok {
		b.structured["bizid"] = config.ServiceConfig.DefaultBizID
	}

	if config.ServiceConfig.IsLocalBizLog {
		// map to fields
		fields := log.Fields{}
		for key, value := range b.structured {
			fields[key] = value
		}
		log.WithFields(fields).Info()
	}
	topic := b.structured["logtype"]
	value := buildStrLog(b.structured)

	header := map[string]string{
		"collectionName": b.logType,
	}
	headerStr := buildStrLog(header)

	logValue := fmt.Sprintf("%010d%010d%s%s", len(headerStr), len(value), headerStr, value)
	// send to kafka
	if config.ServiceConfig.IsRemoteBizLog && producer != nil {
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(logValue)},
			chain,
		)
		if err != nil {
			log.Errorf("kafka produce failed, err: %v", err)
		}
	}
}

func buildStrLog(l map[string]string) string {
	logStr := ""
	for k, v := range l {
		logStr += trimStr(k) + "~" + trimStr(v) + string([]byte{31})
	}
	return logStr
}

func trimStr(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, string([]byte{31}), "")
	return str
}
