package config

import (
	"github.com/opensourceways/message-push/common/kafka"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
)

var GiteeConfigInstance GiteeConfig

type GiteeConfig struct {
	Kafka            kafka.ConsumeConfig      `json:"kafka"`
	HWCloudMsgConfig pushSdk.HWCloudMsgConfig `json:"hw_cloud_msg"`
	EmailConfig      pushSdk.EmailConfig      `json:"email"`
}

func InitGiteeConfig() {
	cfg := new(GiteeConfig)
	if err := utils.LoadFromYaml("config/gitee_conf.yaml", cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return
	}
	GiteeConfigInstance = *cfg
}
