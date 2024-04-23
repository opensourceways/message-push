package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"message-push/common/kafka"
	"message-push/common/postgresql"
	"message-push/models/dto"
)

// Handler
type Handler interface {
	handle(message []byte) error
}

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
		var raw dto.EurBuildMessageRaw
		msgBodyErr := json.Unmarshal(message.Value, &raw)
		if msgBodyErr != nil {
			return msgBodyErr
		}
		fmt.Println(raw)
		eurBuildEvent := raw.ToCloudEvent()
		kafkaSendErr := kafka.SendMsg("eur_build_event", &eurBuildEvent)
		if kafkaSendErr != nil {
			return kafkaSendErr
		}
		save(eurBuildEvent)
		session.MarkMessage(message, "")
	}
	return nil
}

func save(event dto.EurBuildEvent) {
	do := event.ToCloudEventDO()
	if postgresql.DB().Model(&do).Where("source=?", do.Source, "event_id = ?", do.EventId).Updates(&do).RowsAffected == 0 {
		postgresql.DB().Create(&do)
	}
}
