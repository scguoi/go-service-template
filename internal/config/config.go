package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const configFile = "conf/services.yaml"

type ServiceYaml struct {
	GRPCPort     int    `yaml:"GrpcPort"`
	HTTPPort     int    `yaml:"HttpPort"`
	WSPort       int    `yaml:"WSPort"`
	LogLevel     string `yaml:"LogLevel"`
	LogOut       string `yaml:"LogOut"`
	LogFile      string `yaml:"LogFile"`
	LogMaxDays   int    `yaml:"LogMaxDays"`
	APIPort      int    `yaml:"APIPort"`
	MetricPort   int    `yaml:"MetricPort"`
	AgentOpsPort int    `yaml:"AgentOpsPort"`

	CpName  string `yaml:"CpName"`
	LogType string `yaml:"LogType"`

	IsRemoteBizLog bool                   `yaml:"IsRemoteBizLog"`
	AsyncChainSize int                    `yaml:"AsyncChainSize"`
	KafkaClientExt map[string]interface{} `yaml:"KafkaClientExt"`
}

var ServiceConfig *ServiceYaml

func init() {
	ServiceConfig = new(ServiceYaml)
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("read config file "+configFile+" failed. ", err)
	}
	err = yaml.Unmarshal(fileBytes, ServiceConfig)
	if err != nil {
		log.Fatal("unmarshal service config failed. ", err)
	}
}
