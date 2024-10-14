package config

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
)

var ForumConfigInstance ForumConfig

type ForumConfig struct {
	Kafka kafka.ConsumeConfig `yaml:"kafka"`
	Push  PushConfig          `yaml:"push"`
}

func InitForumConfig(configFile string) {
	cfg := new(ForumConfig)
	if err := utils.LoadFromYaml(configFile, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	ForumConfigInstance = *cfg
}
