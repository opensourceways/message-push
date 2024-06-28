package config

import (
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
)

var EurBuildConfigInstance EurBuildConfig

type EurBuildConfig struct {
	Kafka kafka.ConsumeConfig `yaml:"kafka"`
	Push  PushConfig          `yaml:"push"`
}

func InitEurBuildConfig() {
	cfg := new(EurBuildConfig)
	if err := utils.LoadFromYaml("config/eur_build_conf.yaml", cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	logrus.Info("读取eur配置成功", cfg)
	EurBuildConfigInstance = *cfg
}
