package service

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"message-push/service/push"
	"message-push/service/transfer"
)

func SubscribeEurRaw() {
	_ = kfklib.Subscribe("ssp_test", transfer.Handle, []string{"eur_build_raw"})
}

func SubscribeEurEvent() {
	_ = kfklib.Subscribe("ssp_test", push.Handle, []string{"eur_build_event"})
}
