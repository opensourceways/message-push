package service

import (
	"github.com/IBM/sarama"
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"message-push/common/kafka"
	"message-push/service/push"
	"message-push/service/transfer"
)

//func SubscribeEurRaw() {
//	_ = kfklib.Subscribe("ssp_test", transfer.Handle, []string{"eur_build_raw"})
//}

func SubscribeEurEvent() {
	_ = kfklib.Subscribe("ssp_test", push.Handle, []string{"eur_build_event"})
}

func SubscribeEurRaw() {
	cfg := kafka.ConsumeConfig{
		Topic:   "eur_build_raw",
		Address: "127.0.0.1:9092",
		Group:   "ssp_test",
		Offset:  sarama.OffsetOldest,
	}

	h := transfer.EurGroupHandler{}
	kafka.ConsumeGroup(cfg, &h)
}
