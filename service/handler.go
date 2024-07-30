package service

import (
	"github.com/goccy/go-json"
	"github.com/gocql/gocql"
	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
	"strings"
	"time"
)

func GiteeHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	res := handle(event, config.GiteeConfigInstance.Push)
	return res
}

func EurBuildHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	res := handle(event, config.EurBuildConfigInstance.Push)
	return res
}

func OpenEulerMeetingHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	res := handle(event, config.MeetingConfigInstance.Push)
	return res
}

func handleRelatedUsers(event dto.CloudEvents) {
	raw := make(dto.Raw)
	raw.FromJson(event.Data())
	relatedUsers := strings.Split(event.Extensions()["relatedusers"].(string), ",")
	event.SendInnerMessageByRelatedUsers(relatedUsers)
}
func handle(event dto.CloudEvents, push config.PushConfig) error {
	handleRelatedUsers(event)
	return handleSubcribe(event, push)
}
func handleSubcribe(event dto.CloudEvents, push config.PushConfig) error {
	raw := make(dto.Raw)
	raw.FromJson(event.Data())
	flatRaw := raw.Flatten()
	recipients := event.GetRecipient()
	if recipients == nil || len(recipients) == 0 {
		return nil
	}
	stream.Of(recipients...).Filter(
		func(recipient bo.RecipientConfig) bool {
			if recipient.ModeFilter == nil {
				return true
			}
			return flatRaw.ModeFilter(recipient.ModeFilter)
		},
	).ForEach(func(recipient bo.RecipientConfig) {
		if recipient.NeedMessage {
			res := sendHWCloudMessage(raw, recipient, push.MsgConfig)
			insertData(event, flatRaw, res)
			logrus.Info("send message ", event.ID()+" success", recipient.Message)
		}
		if recipient.NeedMail {
			res := sendMail(event, recipient, push.EmailConfig)
			insertData(event, flatRaw, res)
			logrus.Info("send mail ", event.ID()+" success，接收人", recipient.Mail)
		}
		if recipient.NeedInnerMessage {
			res := sendInnerMessage(event, recipient)
			logrus.Info("send inner message ", event.ID()+" success")
			insertData(event, flatRaw, res)
		}
	})
	return nil
}

func sendHWCloudMessage(raw dto.Raw, recipient bo.RecipientConfig, messageConfig pushSdk.MsgConfig) dto.PushResult {
	templateParas := raw.ToMessageArgs(recipient.MessageTemplate)
	return pushSdk.SendHWCloudMessage(messageConfig, templateParas, recipient)
}

func sendInnerMessage(event dto.CloudEvents, recipient bo.RecipientConfig) dto.PushResult {
	return event.SendInnerMessage(recipient)
}

func sendMail(event dto.CloudEvents, recipient bo.RecipientConfig, emailConfig pushSdk.EmailConfig) dto.PushResult {
	return pushSdk.SendEmail(event.Extensions()["title"].(string),
		event.Extensions()["summary"].(string), recipient, emailConfig)
}

func insertData(eurBuildEvent dto.CloudEvents, flatRaw dto.FlatRaw, result dto.PushResult) {
	stringifyMap := flatRaw.StringifyMap()
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
