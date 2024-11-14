package service

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-push/config"
)

func SubscribeEurEvent() {
	logrus.Info("subscribing to eur topic")
	_ = kfklib.Subscribe(config.EurBuildConfigInstance.Kafka.Group, EurBuildHandle, []string{config.EurBuildConfigInstance.Kafka.Topic})
}

func SubscribeGiteeEvent() {
	logrus.Info("subscribing to gitee topic")
	_ = kfklib.Subscribe(config.GiteeConfigInstance.Kafka.Issue.Group, GiteeHandle, []string{config.GiteeConfigInstance.Kafka.Issue.Topic, config.GiteeConfigInstance.Kafka.Push.Topic, config.GiteeConfigInstance.Kafka.PR.Topic, config.GiteeConfigInstance.Kafka.Note.Topic})
}

func SubscribeMeetingEvent() {
	logrus.Info("subscribing to meeting topic")
	_ = kfklib.Subscribe(config.MeetingConfigInstance.Kafka.Group, OpenEulerMeetingHandle, []string{config.MeetingConfigInstance.Kafka.Topic})
}

func SubscribeCVEEvent() {
	logrus.Info("subscribing to cve topic")
	_ = kfklib.Subscribe(config.CVEConfigInstance.Kafka.Group, CVEHandle, []string{config.CVEConfigInstance.Kafka.Topic})
}

func SubscribeForumEvent() {
	logrus.Info("subscribing to forum topic")
	_ = kfklib.Subscribe(config.ForumConfigInstance.Kafka.Group, ForumHandle, []string{config.ForumConfigInstance.Kafka.Topic})
}
