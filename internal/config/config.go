package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var _configFile = "conf/services.yaml"

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

	IsRemoteBizLog bool `yaml:"IsRemoteBizLog"`
	IsLocalBizLog  bool `yaml:"IsLocalBizLog"`

	AsyncChainSize int                    `yaml:"AsyncChainSize"`
	KafkaClientExt map[string]interface{} `yaml:"KafkaClientExt"`

	DefaultAppID       string `yaml:"DefaultAppID"`
	DefaultProductLine string `yaml:"DefaultProductLine"`
	DefaultBizID       string `yaml:"DefaultBizID"`
}

var ServiceConfig *ServiceYaml

func Load() {
	loadConfig(_configFile)
}

func LoadWithYaml(content []byte) {
	ServiceConfig = new(ServiceYaml)
	err := yaml.Unmarshal(content, ServiceConfig)
	if err != nil {
		log.Fatal("unmarshal service config failed. ", err)
	}
}

func loadConfig(file string) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("read config file "+file+" failed. ", err)
	}
	LoadWithYaml(fileBytes)
}
