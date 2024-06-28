package config

import (
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
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

func InitGiteeConfig() {
	cfg := new(GiteeConfig)
	if err := utils.LoadFromYaml("config/gitee_conf.yaml", cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	logrus.Info("读取gitee配置成功", cfg)
	GiteeConfigInstance = *cfg
}
