package service

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/message-push/config"
	"github.com/sirupsen/logrus"
)

func SubscribeEurEvent() {
	logrus.Info("subscribing to eur topic")
	_ = kfklib.Subscribe(config.EurBuildConfigInstance.Kafka.Group, EurBuildHandle, []string{config.EurBuildConfigInstance.Kafka.Topic})
}

func SubscribeGiteeEvent() {
	logrus.Info("subscribing to gitee topic")
	_ = kfklib.Subscribe(config.GiteeConfigInstance.Kafka.Issue.Group, GiteeHandle, []string{config.GiteeConfigInstance.Kafka.Issue.Topic, config.GiteeConfigInstance.Kafka.Push.Topic, config.GiteeConfigInstance.Kafka.PR.Topic, config.GiteeConfigInstance.Kafka.Note.Topic})
}
