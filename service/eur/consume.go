package eur

import (
	"github.com/IBM/sarama"
	"message-push/models/messageadapter"
	"message-push/service/hanlder"
)

func ConsumeEur() {
	cfg := messageadapter.ConsumeConfig{
		Topic:   "eur_build_event",
		Address: "0.0.0.0:9092",
		Group:   "ssp_test",
		Offset:  sarama.OffsetOldest,
	}

	h := hanlder.EurGroupHandler{}
	messageadapter.ConsumeGroup(cfg, &h)
}
