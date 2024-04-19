package push

import (
	"context"
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
	"message-push/common/cassandra"
	"message-push/common/pushSdk"
	"message-push/models/bo"
	"message-push/models/dto"
	"message-push/utils"
	"strconv"
	"time"
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
	flatRaw := eurBuildRaw.Flatten()
	stream.Of(subscribes...).Filter(
		func(item bo.SubscribePushConfig) bool {
			return utils.ModeFilter(flatRaw, item.ModeFilter)
		},
	).ForEach(func(item bo.SubscribePushConfig) {
		var cfg []bo.PushConfig
		_ = json.Unmarshal(item.PushConfigs, &cfg)
		for _, push := range cfg {
			switch push.PushType {
			case "phone":
				context.TODO()
			case "message":
				//sendHWCloudMessage(eurBuildRaw, push)
				insertData(event, flatRaw, push)
			case "api":
				context.TODO()
			default:
				logrus.Info("不支持的推送类型:", push.PushType)
			}
		}
	})
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

func insertData(eurBuildEvent dto.EurBuildEvent, flatRaw map[string]interface{}, push bo.PushConfig) {
	stringifyMap := utils.StringifyMap(flatRaw)
	insert := `INSERT INTO message_center.message_push_record
(id, created_at, data, event_id, push_address, push_state, push_time,
 push_type, source)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`
	err := cassandra.Session().
		Query(
			insert,
			gocql.TimeUUID(),
			time.Now(),
			stringifyMap,
			eurBuildEvent.ID(),
			push.PushAddress,
			"success",
			time.Now(),
			push.PushType,
			eurBuildEvent.Source(),
		).
		Exec()
	if err != nil {
		panic(nil)
		return

	}
}
