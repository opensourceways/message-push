package service

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/service/push"
)

func SubscribeEurEvent() {
	config.InitEurBuildConfig()
	_ = kfklib.Subscribe(config.EurBuildConfigInstance.Kafka.Group, push.Handle, []string{config.EurBuildConfigInstance.Kafka.Topic})
}
