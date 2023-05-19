package logc

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"template/internal/config"
	"time"
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

func (b *BizLog) LoggerStructured(key string, value string) {
	b.structured[key] = value
}

func (b *BizLog) LoggerStructuredBatch(logs map[string]string) {
	for k, v := range logs {
		b.structured[k] = v
	}
}

func (b *BizLog) LoggerEnd() {
	raw, err := json.Marshal(b.raw)
	if err != nil {
		log.Errorf("json marshal failed, err: %v", err)
	}
	b.structured["inner"] = string(raw)
	b.etm = time.Now()
	b.utm = int32(b.etm.Sub(b.stm).Nanoseconds() / 1e6)
	b.structured["etm"] = b.etm.Format("2006-01-02 15:04:05.000")
	b.structured["utm"] = string(b.utm)
	// map to fields
	fields := log.Fields{}
	for key, value := range b.structured {
		fields[key] = value
	}
	log.WithFields(fields).Info()
}
