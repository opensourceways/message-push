package push

import (
	"context"
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
	"strconv"
	"strings"
	"time"
)

func Handle(payload []byte, _ map[string]string) error {
	eurBuildEvent := dto.NewCloudEvents()
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
	if subscribes == nil || len(subscribes) == 0 {
		return
	}
	flatRaw := eurBuildRaw.Flatten()
	stream.Of(subscribes...).Filter(
		func(item bo.SubscribePushConfig) bool {
			if item.ModeFilter == nil {
				return true
			}
			return utils.ModeFilter(flatRaw, item.ModeFilter)
		},
	).ForEach(func(item bo.SubscribePushConfig) {
		var cfg []bo.PushConfig
		_ = json.Unmarshal(item.PushConfigs, &cfg)
		for _, push := range cfg {
			switch push.PushType {
			case "inner_message":
				res := sendInnerMessage(event, item)
				insertData(event, flatRaw, push, item.RecipientId, res)
			case "phone":
				context.TODO()
			case "message":
				res := sendHWCloudMessage(eurBuildRaw, push)
				insertData(event, flatRaw, push, item.RecipientId, res)
			case "api":
				context.TODO()
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
}

func sendInnerMessage(eurBuildEvent dto.CloudEvents, config bo.SubscribePushConfig) dto.PushResult {
	return eurBuildEvent.SendInnerMessage(config.RecipientId)
}

func insertData(eurBuildEvent dto.CloudEvents, flatRaw map[string]interface{}, push bo.PushConfig, recipient string, result dto.PushResult) {
	stringifyMap := utils.StringifyMap(flatRaw)
	insert := `insert into message_push_record (recipient_id, time_uuid, created_at, event_data, event_data_content_type,
                                 event_data_schema, event_id, event_source, event_source_url, event_spec_version,
                                 event_time, event_type, event_user, push_address, push_state, push_time, push_type,
                                 remark, title, summary)
values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
`
	err := cassandra.Session().
		Query(
			insert,
			recipient,
			gocql.TimeUUID(),
			time.Now(),
			stringifyMap,
			eurBuildEvent.DataContentType(),
			eurBuildEvent.DataSchema(),
			eurBuildEvent.ID(),
			eurBuildEvent.Source(),
			eurBuildEvent.Extensions()["sourceurl"].(string),
			eurBuildEvent.SpecVersion(),
			eurBuildEvent.Time(),
			eurBuildEvent.Type(),
			eurBuildEvent.Extensions()["user"].(string),
			push.PushAddress,
			result.Res,
			result.Time,
			push.PushType,
			result.Remark,
			eurBuildEvent.Extensions()["title"].(string),
			eurBuildEvent.Extensions()["summary"].(string),
		).
		Exec()
	if err != nil {
		panic(nil)
		return
	}
}
