package config

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
)

var GiteeConfigInstance GiteeConfig

type GiteeConfig struct {
	Kafka struct {
		Issue kafka.ConsumeConfig `yaml:"issue"`
		Push  kafka.ConsumeConfig `yaml:"push"`
		PR    kafka.ConsumeConfig `yaml:"pr"`
		Note  kafka.ConsumeConfig `yaml:"note"`
	} `yaml:"kafka"`
	Push PushConfig `yaml:"push"`
}

func InitGiteeConfig(configFile string) {
	cfg := new(GiteeConfig)
	if err := utils.LoadFromYaml(configFile, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	GiteeConfigInstance = *cfg
}
