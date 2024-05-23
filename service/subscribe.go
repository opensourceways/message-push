package service

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/message-push/service/push"
)

func SubscribeEurEvent() {
	_ = kfklib.Subscribe("ssp_test", push.Handle, []string{"eur_build_event"})
}
