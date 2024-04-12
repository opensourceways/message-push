package hanlder

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/opensourceways/kafka-lib/mq"
	"message-push/common/postgresql"
	"message-push/models/event"
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
		var eurBuildEvent event.EurBuildEvent

		msgBodyErr := json.Unmarshal(msg.Body, &eurBuildEvent)
		if msgBodyErr != nil {
			return err
		}
		fmt.Printf("Received message with offset %d: %s\n", message.Offset, eurBuildEvent)

		transferErr := publishEurEvent(eurBuildEvent)
		if transferErr != nil {
			return transferErr
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func publishEurEvent(event event.EurBuildEvent) error {
	messageadapter.SendHWCloudMessage(&event)
	return nil
}

func save(raw event.EurBuildRaw) {
	do := raw.ToCloudEventDO()
	res := postgresql.DB().Table("message_center.cloud_event_message").Create(&do)
	fmt.Println(res)
}
