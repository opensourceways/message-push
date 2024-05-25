package config

import (
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
)

var EurBuildConfigInstance EurBuildConfig

type EurBuildConfig struct {
	Kafka            kafka.ConsumeConfig      `json:"kafka"`
	HWCloudMsgConfig pushSdk.HWCloudMsgConfig `json:"hw_cloud_msg"`
}

func InitEurBuildConfig() {
	cfg := new(EurBuildConfig)
	if err := utils.LoadFromYaml("config/eur_build_conf.yaml", cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	EurBuildConfigInstance = *cfg
}
