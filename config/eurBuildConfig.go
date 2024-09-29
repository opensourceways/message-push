package config

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
)

var EurBuildConfigInstance EurBuildConfig

type EurBuildConfig struct {
	Kafka kafka.ConsumeConfig `yaml:"kafka"`
	Push  PushConfig          `yaml:"push"`
}

func InitEurBuildConfig(configFile string) {
	cfg := new(EurBuildConfig)
	if err := utils.LoadFromYaml(configFile, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	EurBuildConfigInstance = *cfg
}
