package push

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"message-push/common/pushSdk"
	"message-push/models/bo"
	"message-push/models/dto"
	"strconv"
)

func Handle(payload []byte, _ map[string]string) error {
	var eurBuildEvent dto.EurBuildEvent
	msgBodyErr := json.Unmarshal(payload, &eurBuildEvent)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	publishMessage(eurBuildEvent)
	return nil
}

func publishMessage(event dto.EurBuildEvent) {
	var eurBuildRaw dto.EurBuildRaw
	_ = json.Unmarshal(event.Data(), &eurBuildRaw)
	subscribes := event.GetSubscribe()
	for _, subscribe := range subscribes {
		var cfg []bo.PushConfig
		_ = json.Unmarshal(subscribe.PushConfigs, &cfg)
		for _, push := range cfg {
			switch push.PushType {
			case "phone":
				context.TODO()
			case "message":
				sendHWCloudMessage(eurBuildRaw, push)
			case "api":
				context.TODO()
			default:
				logrus.Info("不支持的推送类型:", push.PushType)
			}
		}
	}
}

func sendHWCloudMessage(eurBuildRaw dto.EurBuildRaw, push bo.PushConfig) {
	masConfig := pushSdk.NewTestConfig()
	templateParas := []string{
		strconv.Itoa(eurBuildRaw.Body.Build),
		"success",
		eurBuildRaw.Body.Owner,
		eurBuildRaw.Body.Copr,
		strconv.Itoa(eurBuildRaw.Body.Build),
	}
	pushSdk.SendHWCloudMessage(masConfig, templateParas, push.PushAddress)
}
