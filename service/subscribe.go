package service

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/message-push/config"
)

func SubscribeEurEvent() {
	_ = kfklib.Subscribe(config.EurBuildConfigInstance.Kafka.Group, EurBuildHandle, []string{config.EurBuildConfigInstance.Kafka.Topic})
}

func SubscribeGiteeIssueEvent() {
	_ = kfklib.Subscribe(config.GiteeConfigInstance.Kafka.Issue.Group, GiteeHandle, []string{config.GiteeConfigInstance.Kafka.Issue.Topic, config.GiteeConfigInstance.Kafka.Push.Topic})
}
