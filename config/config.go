package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const configFile = "conf/services.yaml"

type ServiceYaml struct {
	GRPCPort int    `yaml:"GrpcPort"`
	HTTPPort int    `yaml:"HttpPort"`
	LogLevel string `yaml:"LogLevel"`
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
