package hanlder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
	"message-push/models/bo"
	"message-push/models/dto"
	"message-push/models/messageadapter"
)

type EurHandler struct{}

func (eurHandler *EurHandler) handle(message []byte) error {
	fmt.Println(message)
	return nil
}

type EurGroupHandler struct{}

func (h EurGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h EurGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h EurGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

	json.Unmarshal(event.Data(), &eurBuildRaw)
	subscribes := event.GetSubscribe()
	stream.Of(subscribes...).Filter(
		func(item bo.SubscribePushConfig) bool {
			return eurBuildRaw.ModeFilter(item.ModeFilter)
		},
	).ForEach(
		func(subscribe bo.SubscribePushConfig) {
			var cfg []bo.PushConfig
			json.Unmarshal(subscribe.PushConfigs, &cfg)
			for _, push := range cfg {
				switch push.PushType {
				case "phone":
					context.TODO()
				case "message":
					messageadapter.SendHWCloudMessage(&eurBuildRaw, push.PushAddress)
				case "api":
					context.TODO()
				default:
					logrus.Info("不支持的推送类型:", push.PushType)
				}
			}
		},
	)
}
