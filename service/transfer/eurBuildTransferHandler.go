package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/opensourceways/kafka-lib/mq"
	"message-push/common/kafka"
	"message-push/common/postgresql"
	"message-push/models/dto"
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
		var raw dto.EurBuildRaw

		msgBodyErr := json.Unmarshal(msg.Body, &raw)
		if msgBodyErr != nil {
			return err
		}
		fmt.Printf("Received message with offset %d: %s\n", message.Offset, raw)

		transferErr := publishEurEvent(raw)
		if transferErr != nil {
			return transferErr
		}
		save(raw)
		session.MarkMessage(message, "")
	}
	return nil
}

func publishEurEvent(raw dto.EurBuildRaw) error {
	eurBuildEvent := raw.ToCloudEvent()
	sendErr := kafka.SendMsg("eur_build_event", &eurBuildEvent)
	if sendErr != nil {
		return sendErr
	}
	return nil
}

func save(raw dto.EurBuildRaw) {
	do := raw.ToCloudEventDO()
	if postgresql.DB().Model(&do).Where("event_id = ?", do.EventId).Updates(&do).RowsAffected == 0 {
		postgresql.DB().Create(&do)
	}
}
