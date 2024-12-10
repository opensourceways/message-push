package config

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/utils"
)

var MeetingConfigInstance MeetingConfig

type MeetingConfig struct {
	Kafka kafka.ConsumeConfig `yaml:"kafka"`
	Push  PushConfig          `yaml:"push"`
}

func InitMeetingConfig(configFile string) {
	cfg := new(MeetingConfig)
	logrus.Infof("configFile is %v", configFile)
	if err := utils.LoadFromYaml(configFile, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	logrus.Infof("the cfg is %v", cfg)
	MeetingConfigInstance = *cfg
}
