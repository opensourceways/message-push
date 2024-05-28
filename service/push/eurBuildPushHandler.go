package push

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
	"strconv"
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
	recipients := event.GetRecipient()
	if recipients == nil || len(recipients) == 0 {
		return
	}
	flatRaw := eurBuildRaw.Flatten()
	stream.Of(recipients...).Filter(
		func(recipient bo.RecipientConfig) bool {
			if recipient.ModeFilter == nil {
				return true
			}
			return utils.ModeFilter(flatRaw, recipient.ModeFilter)
		},
	).ForEach(func(recipient bo.RecipientConfig) {
		if recipient.NeedMessage {
			res := sendHWCloudMessage(eurBuildRaw, recipient)
			insertData(event, flatRaw, res)
			logrus.Info("send message ", event.ID()+" success")
		}
		if recipient.NeedInnerMessage {
			res := sendInnerMessage(event, recipient)
			logrus.Info("send inner message ", event.ID()+" success")
			insertData(event, flatRaw, res)
		}
	})
}

func sendHWCloudMessage(eurBuildRaw dto.EurBuildMessageRaw, recipient bo.RecipientConfig) dto.PushResult {
	status := ""
	switch eurBuildRaw.Body.Status {
	case 0:
		status = "失败"
	case 1:
		status = "成功"
	case 3:
		status = "开始"
	default:
		status = "未知"
	}
	templateParas := []string{
		strconv.Itoa(eurBuildRaw.Body.Build),
		status,
		eurBuildRaw.Body.Owner,
		eurBuildRaw.Body.Copr,
		strconv.Itoa(eurBuildRaw.Body.Build),
	}
	return pushSdk.SendHWCloudMessage(config.EurBuildConfigInstance.HWCloudMsgConfig, templateParas, recipient)
}

func sendInnerMessage(eurBuildEvent dto.CloudEvents, recipient bo.RecipientConfig) dto.PushResult {
	return eurBuildEvent.SendInnerMessage(recipient)
}

func insertData(eurBuildEvent dto.CloudEvents, flatRaw map[string]interface{}, result dto.PushResult) {
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
			result.RecipientId,
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
			result.PushAddress,
			result.Res,
			result.Time,
			result.PushType,
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
