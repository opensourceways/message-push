package push

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
	"message-push/common/pushSdk"
	"message-push/models/bo"
	"message-push/models/dto"
	"strconv"
)

type EurBuildPushHandler struct{}

func (h EurBuildPushHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h EurBuildPushHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h EurBuildPushHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg mq.Message
		err := json.Unmarshal(message.Value, &msg)
		if err != nil {
			return err
		}
		var eurBuildEvent dto.EurBuildEvent
		msgBodyErr := json.Unmarshal(msg.Body, &eurBuildEvent)
		if msgBodyErr != nil {
			return err
		}
		fmt.Printf("Received message with offset %d: %s\n", message.Offset, eurBuildEvent)
		publishEurEvent(eurBuildEvent)
		session.MarkMessage(message, "")
	}
	return nil
}

func publishEurEvent(event dto.EurBuildEvent) {
	var eurBuildRaw dto.EurBuildRaw
	_ = json.Unmarshal(event.Data(), &eurBuildRaw)
	subscribes := event.GetSubscribe()
	stream.Of(subscribes...).Filter(
		func(item bo.SubscribePushConfig) bool {
			return eurBuildRaw.ModeFilter(item.ModeFilter)
		},
	).ForEach(
		func(subscribe bo.SubscribePushConfig) {
			var cfg []bo.PushConfig
			_ = json.Unmarshal(subscribe.PushConfigs, &cfg)
			for _, push := range cfg {
				switch push.PushType {
				case "phone":
					context.TODO()
				case "message":
					sendHWCloudMessage(eurBuildRaw, push)
				case "api":
					context.TODO()
				default:
					logrus.Info("不支持的推送类型:", push.PushType)
				}
			}
		},
	)
}

func sendHWCloudMessage(eurBuildRaw dto.EurBuildRaw, push bo.PushConfig) {
	masConfig := pushSdk.NewTestConfig()
	templateParas := []string{
		strconv.Itoa(eurBuildRaw.Body.Build),
		"success",
		eurBuildRaw.Body.Owner,
		eurBuildRaw.Body.Copr,
		strconv.Itoa(eurBuildRaw.Body.Build),
	}
	pushSdk.SendHWCloudMessage(masConfig, templateParas, push.PushAddress)
}
