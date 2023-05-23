package logc

import (
	"strings"
	"template/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	config.LoadWithYaml([]byte(`
GrpcPort: 1234
HttpPort: 8080
WSPort: 5678
LogLevel: "debug"
LogOut: "stdout"
LogFile: "log.txt"
LogMaxDays: 7
APIPort: 4321
MetricPort: 56789
AgentOpsPort: 6789
CpName: "test"
LogType: "json"
IsRemoteBizLog: false
IsLocalBizLog: false
AsyncChainSize: 10
KafkaClientExt:
  bootstrap.servers: "localhost:9092"
DefaultAppID: "test-app"
DefaultProductLine: "test-product"
DefaultBizID: "test-biz"
`))
	m.Run()
}

func TestNewBizLog(t *testing.T) {
	b := NewBizLog()
	assert.NotNil(t, b)
	assert.Equal(t, config.ServiceConfig.LogType, b.logType)
	assert.Equal(t, config.ServiceConfig.CpName, b.cpName)
	assert.NotEqual(t, time.Time{}, b.stm)
	assert.Empty(t, b.raw)
	assert.NotEmpty(t, b.structured["logtype"])
	assert.NotEmpty(t, b.structured["cpname"])
	assert.NotEmpty(t, b.structured["stm"])
}

func TestBizLog_LoggerRaw(t *testing.T) {
	b := NewBizLog()
	b.LoggerRaw("test raw log")
	assert.Equal(t, []string{"test raw log"}, b.raw)
}

func TestBizLog_Printf(t *testing.T) {
	b := NewBizLog()
	b.Printf("test %s log", "formatted")
	assert.Equal(t, []string{"test formatted log"}, b.raw)
}

func TestBizLog_Print(t *testing.T) {
	b := NewBizLog()
	b.Print("test", "print", "log")
	assert.Equal(t, []string{"test print log"}, b.raw)
}

func TestBizLog_LoggerStructured(t *testing.T) {
	b := NewBizLog()
	b.LoggerStructured("key", "value")
	assert.Equal(t, "value", b.structured["key"])
}

func TestBizLog_LoggerStructuredBatch(t *testing.T) {
	b := NewBizLog()
	b.LoggerStructuredBatch(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	assert.Equal(t, "value1", b.structured["key1"])
	assert.Equal(t, "value2", b.structured["key2"])
}

func TestBizLog_LoggerEnd(t *testing.T) {
	b := NewBizLog()
	b.LoggerEnd()
	assert.NotEmpty(t, b.structured["etm"])
	assert.NotEmpty(t, b.structured["utm"])
	assert.NotEmpty(t, b.structured["traceid"])
	assert.NotEmpty(t, b.structured["appid"])
	assert.NotEmpty(t, b.structured["productline"])
	assert.NotEmpty(t, b.structured["bizid"])
}

func TestTrimStr(t *testing.T) {
	assert.Equal(t, "test", trimStr("test"))
	assert.Equal(t, "test", trimStr("test\n"))
	assert.Equal(t, "test", trimStr("test\r"))
	assert.Equal(t, "test", trimStr("test\x1f"))
}

func BenchmarkBizLog_LoggerEnd2K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bz := NewBizLog()
		bz.LoggerStructured("key1", strings.Repeat("value1", 100))
		bz.LoggerStructured("key2", strings.Repeat("value2", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerEnd()
	}
}

func BenchmarkBizLog_LoggerEnd1K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bz := NewBizLog()
		bz.LoggerStructured("key1", strings.Repeat("value1", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerEnd()
	}
}

func BenchmarkBizLog_LoggerEnd10K(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bz := NewBizLog()
		bz.LoggerStructured("key1", strings.Repeat("value1", 100))
		bz.LoggerStructured("key2", strings.Repeat("value2", 100))
		bz.LoggerStructured("key3", strings.Repeat("value3", 100))
		bz.LoggerStructured("key4", strings.Repeat("value4", 100))
		bz.LoggerStructured("key5", strings.Repeat("value5", 100))
		bz.LoggerStructured("key6", strings.Repeat("value6", 100))
		bz.LoggerStructured("key7", strings.Repeat("value7", 100))
		bz.LoggerStructured("key8", strings.Repeat("value8", 100))
		bz.LoggerStructured("key9", strings.Repeat("value9", 100))
		bz.LoggerStructured("key10", strings.Repeat("value10", 100))
		bz.LoggerStructured("key11", strings.Repeat("value11", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerRaw(strings.Repeat("raw", 100))
		bz.LoggerEnd()
	}
}
