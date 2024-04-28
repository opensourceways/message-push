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
	"strings"
	"time"
)

func Handle(payload []byte, _ map[string]string) error {
	var eurBuildEvent dto.CloudEvents
	msgBodyErr := json.Unmarshal(payload, &eurBuildEvent)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	publishMessage(eurBuildEvent)
	return nil
}

func publishMessage(event dto.CloudEvents) {
	var eurBuildRaw dto.EurBuildMessageRaw
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
				res := sendHWCloudMessage(eurBuildRaw, push)
				insertData(event, flatRaw, push, item.RecipientId, res)
			case "api":
				res := dto.PushResult{Res: dto.Succeed}
				insertData(event, flatRaw, push, item.RecipientId, res)
			default:
				logrus.Info("不支持的推送类型:", push.PushType)
			}
		}
	})
}

func sendHWCloudMessage(eurBuildRaw dto.EurBuildMessageRaw, push bo.PushConfig) dto.PushResult {
	masConfig := pushSdk.NewTestConfig()
	topicArray := strings.Split(eurBuildRaw.Topic, ".")
	templateParas := []string{
		strconv.Itoa(eurBuildRaw.Body.Build),
		topicArray[len(topicArray)-1],
		eurBuildRaw.Body.Owner,
		eurBuildRaw.Body.Copr,
		strconv.Itoa(eurBuildRaw.Body.Build),
	}
	return pushSdk.SendHWCloudMessage(masConfig, templateParas, push.PushAddress)
	return dto.PushResult{
		Res: dto.Succeed,
	}
}

func insertData(eurBuildEvent dto.CloudEvents, flatRaw map[string]interface{}, push bo.PushConfig, recipient string, result dto.PushResult) {
	stringifyMap := utils.StringifyMap(flatRaw)
	insert := `insert into message_push_record
			   (
			    recipient_id,
 				source,
 				time_uuid,
 				created_at,
 				data, 
			    event_id,
 				push_address,
 				push_state,
 				push_time,
 				push_type,
			    remark
 				)
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?);
`
	err := cassandra.Session().
		Query(
			insert,
			recipient,
			eurBuildEvent.Source(),
			gocql.TimeUUID(),
			time.Now(),
			stringifyMap,
			eurBuildEvent.ID(),
			push.PushAddress,
			result.Res,
			time.Now(),
			push.PushType,
			result.Remark,
		).
		Exec()
	if err != nil {
		panic(nil)
		return

	}
}
