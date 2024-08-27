package config

import (
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
)

var CVEConfigInstance CVEConfig

type CVEConfig struct {
	Kafka kafka.ConsumeConfig `yaml:"kafka"`
	Push  PushConfig          `yaml:"push"`
}

func InitCVEConfig(configFile string) {
	cfg := new(CVEConfig)
	if err := utils.LoadFromYaml(configFile, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	CVEConfigInstance = *cfg
}
