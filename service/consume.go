package service

import (
	"github.com/IBM/sarama"
	"message-push/common/kafka"
	"message-push/service/push"
	"message-push/service/transfer"
)

func ConsumeEurBuildPush() {
	cfg := kafka.ConsumeConfig{
		Topic:   "eur_build_event",
		Address: "0.0.0.0:9092",
		Group:   "ssp_test",
		Offset:  sarama.OffsetOldest,
	}

	h := push.EurBuildPushHandler{}
	kafka.ConsumeGroup(cfg, &h)
}

func ConsumeEurBuildTransfer() {
	cfg := kafka.ConsumeConfig{
		Topic:   "eur_build_raw",
		Address: "0.0.0.0:9092",
		Group:   "ssp_test",
		Offset:  sarama.OffsetNewest,
	}
	h := transfer.EurBuildPushHandler{}
	kafka.ConsumeGroup(cfg, &h)
}
