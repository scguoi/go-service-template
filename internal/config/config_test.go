package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// create a temporary file for testing
	tmpfile, err := ioutil.TempFile("", "test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// write test data to the temporary file
	_, err = tmpfile.Write([]byte(`
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
IsRemoteBizLog: true
IsLocalBizLog: false
AsyncChainSize: 10
KafkaClientExt:
  bootstrap.servers: "localhost:9092"
DefaultAppID: "test-app"
DefaultProductLine: "test-product"
DefaultBizID: "test-biz"
`))
	if err != nil {
		t.Fatal(err)
	}
	// close the file to flush the contents to disk
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// load the config
	loadConfig(tmpfile.Name())

	// check the values
	assert.Equal(t, 1234, ServiceConfig.GRPCPort)
	assert.Equal(t, 8080, ServiceConfig.HTTPPort)
	assert.Equal(t, 5678, ServiceConfig.WSPort)
	assert.Equal(t, "debug", ServiceConfig.LogLevel)
	assert.Equal(t, "stdout", ServiceConfig.LogOut)
	assert.Equal(t, "log.txt", ServiceConfig.LogFile)
	assert.Equal(t, 7, ServiceConfig.LogMaxDays)
	assert.Equal(t, 4321, ServiceConfig.APIPort)
	assert.Equal(t, 56789, ServiceConfig.MetricPort)
	assert.Equal(t, 6789, ServiceConfig.AgentOpsPort)
	assert.Equal(t, "test", ServiceConfig.CpName)
	assert.Equal(t, "json", ServiceConfig.LogType)
	assert.Equal(t, true, ServiceConfig.IsRemoteBizLog)
	assert.Equal(t, false, ServiceConfig.IsLocalBizLog)
	assert.Equal(t, 10, ServiceConfig.AsyncChainSize)
	assert.Equal(t, "localhost:9092", ServiceConfig.KafkaClientExt["bootstrap.servers"])
	assert.Equal(t, "test-app", ServiceConfig.DefaultAppID)
	assert.Equal(t, "test-product", ServiceConfig.DefaultProductLine)
	assert.Equal(t, "test-biz", ServiceConfig.DefaultBizID)
}
