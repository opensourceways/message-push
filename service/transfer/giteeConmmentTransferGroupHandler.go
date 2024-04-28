package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"message-push/common/kafka"
	"message-push/models/dto"
)

type GiteeHandler struct{}

func (giteeHandler *GiteeHandler) handle(message []byte) error {
	fmt.Println(message)
	return nil
}

type GiteeGroupHandler struct{}

func (h GiteeGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h GiteeGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h GiteeGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var raw dto.Hook
		msgBodyErr := json.Unmarshal(message.Value, &raw)
		if msgBodyErr != nil {
			return msgBodyErr
		}
		fmt.Println(raw)
		event := raw.Comment.ToCloudEventCloudEvents()
		kafkaSendErr := kafka.SendMsg("gitee_comment_event", &event)
		if kafkaSendErr != nil {
			return kafkaSendErr
		}
		event.SaveDb()
		session.MarkMessage(message, "")
	}
	return nil
}
