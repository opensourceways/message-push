package config

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
)

var CertificationConfigInstance CertificationConfig

type CertificationConfig struct {
	Kafka kafka.ConsumeConfig `yaml:"kafka"`
	Push  PushConfig          `yaml:"push"`
}

func InitCertificationConfig(configFile string) {
	cfg := new(CertificationConfig)
	if err := utils.LoadFromYaml(configFile, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	CertificationConfigInstance = *cfg
}
