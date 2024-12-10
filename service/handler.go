package service

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"

	"github.com/opensourceways/message-push/common/cassandra"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
)

func GiteeHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	err := HandleAll(event, config.GiteeConfigInstance.Push)
	if err != nil {
		return err
	}
	return nil
}

func EurBuildHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	err := HandleAll(event, config.EurBuildConfigInstance.Push)
	if err != nil {
		return err
	}
	return nil
}

func OpenEulerMeetingHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	err := HandleAll(event, config.MeetingConfigInstance.Push)
	if err != nil {
		return err
	}
	return nil
}

func CVEHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	err := HandleAll(event, config.CVEConfigInstance.Push)
	if err != nil {
		return err
	}
	return nil
}

func ForumHandle(payload []byte, _ map[string]string) error {
	event := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &event)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	err := HandleAll(event, config.ForumConfigInstance.Push)
	if err != nil {
		return err
	}
	return nil
}

func HandleAll(event dto.CloudEvents, push config.PushConfig) error {
	err := HandleRelated(event)
	if err != nil {
		return err
	}
	err = HandleSubscribe(event, push)
	if err != nil {
		return err
	}
	err = HandleTodo(event)
	if err != nil {
		return err
	}
	err = HandleFollow(event)
	if err != nil {
		return err
	}
	return nil
}

func HandleSubscribe(event dto.CloudEvents, push config.PushConfig) error {
	raw := make(dto.Raw)
	raw.FromJson(event.Data())
	recipients := event.GetSubscribeFromDB()

	if recipients == nil || len(recipients) == 0 {
		return nil
	}
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true, // 启用美化输出
	})
	flatRaw := raw.Flatten()
	processedRecipients := make(map[string]struct{})

	// 遍历接收者
	stream.Of(recipients...).ForEach(func(item bo.RecipientPushConfig) {
		recipientKey := item.RecipientId
		if _, exists := processedRecipients[recipientKey]; !exists {
			logrus.Infof("send email, %v, %v", item.NeedMail, item.Mail)
			isFilter := flatRaw.ModeFilter(item.ModeFilter)
			if isFilter {
				HandleMail(event, flatRaw, item, push)
				if item.NeedMail {
					processedRecipients[recipientKey] = struct{}{}
				}
			}
		}
	})
	return nil
}

func HandleRelated(event dto.CloudEvents) error {
	raw := make(dto.Raw)
	raw.FromJson(event.Data())
	recipients := event.GetRelatedFromDB()

	if recipients == nil || len(recipients) == 0 {
		return nil
	}
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true, // 启用美化输出
	})
	logrus.Info(recipients)
	flatRaw := raw.Flatten()
	processedInnerRecipients := make(map[string]struct{}) // 用于追踪已处理的接收者

	// 遍历接收者
	stream.Of(recipients...).ForEach(func(item bo.RecipientPushConfig) {
		recipientKey := item.RecipientId
		if _, exists := processedInnerRecipients[recipientKey]; !exists {
			HandleInnerMessage(event, flatRaw, item)
			processedInnerRecipients[recipientKey] = struct{}{} // 标记为已处理
		}
	})
	return nil
}

func HandleTodo(event dto.CloudEvents) error {
	raw := make(dto.Raw)
	raw.FromJson(event.Data())
	recipients := event.GetTodoFromDB()
	if recipients == nil || len(recipients) == 0 {
		return nil
	}
	flatRaw := raw.Flatten()
	processedInnerRecipients := make(map[string]struct{}) // 用于追踪已处理的接收者

	// 遍历接收者
	stream.Of(recipients...).ForEach(func(item bo.RecipientPushConfig) {
		recipientKey := item.RecipientId
		if _, exists := processedInnerRecipients[recipientKey]; !exists {
			HandleTodoMessage(event, flatRaw, item)
			processedInnerRecipients[recipientKey] = struct{}{} // 标记为已处理
		}
	})
	return nil
}

func HandleFollow(event dto.CloudEvents) error {
	raw := make(dto.Raw)
	raw.FromJson(event.Data())
	recipients := event.GetFollowFromDB()
	if recipients == nil || len(recipients) == 0 {
		return nil
	}
	flatRaw := raw.Flatten()
	processedInnerRecipients := make(map[string]struct{}) // 用于追踪已处理的接收者

	// 遍历接收者
	stream.Of(recipients...).ForEach(func(item bo.RecipientPushConfig) {
		recipientKey := item.RecipientId
		if _, exists := processedInnerRecipients[recipientKey]; !exists {
			HandleFollowMessage(event, flatRaw, item)
			processedInnerRecipients[recipientKey] = struct{}{} // 标记为已处理
		}
	})
	return nil
}

func HandleInnerMessage(event dto.CloudEvents, flatRaw dto.FlatRaw,
	pushConfig bo.RecipientPushConfig) {
	res := sendInnerMessage(event, pushConfig)
	sendInnerMessageLog := "send inner message %s %s"
	if res.Res == dto.Failed {
		logrus.Info(sendInnerMessageLog, event.ID(), "failed")
	} else {
		logrus.Info(sendInnerMessageLog, event.ID(), "success")
	}
	//insertData(event, flatRaw, res)
}

func HandleTodoMessage(event dto.CloudEvents, flatRaw dto.FlatRaw,
	pushConfig bo.RecipientPushConfig) {
	res := sendTodoMessage(event, pushConfig)
	sendTodoMessageLog := "send inner message %s %s"
	if res.Res == dto.Failed {
		logrus.Info(sendTodoMessageLog, event.ID(), "failed")
	} else {
		logrus.Info(sendTodoMessageLog, event.ID(), "success")
	}
	//insertData(event, flatRaw, res)
}

func HandleFollowMessage(event dto.CloudEvents, flatRaw dto.FlatRaw, pushConfig bo.RecipientPushConfig) {
	res := sendFollowMessage(event, pushConfig)
	sendFollowMessageLog := "send inner message %s %s"
	if res.Res == dto.Failed {
		logrus.Info(sendFollowMessageLog, event.ID(), "failed")
	} else {
		logrus.Info(sendFollowMessageLog, event.ID(), "success")
	}
	//insertData(event, flatRaw, res)
}

func HandleMessage(event dto.CloudEvents, raw dto.Raw, flatRaw dto.FlatRaw,
	pushConfig bo.RecipientPushConfig, push config.PushConfig) {
	if pushConfig.NeedMessage {
		sendMessageLog := "send message %s %s %s"
		res := sendHWCloudMessage(raw, pushConfig, push.MsgConfig)
		if res.Res == dto.Failed {
			logrus.Info(sendMessageLog, event.ID(), "failed", pushConfig.Message)
		} else {
			logrus.Info(sendMessageLog, event.ID(), "success", pushConfig.Message)
		}
		//insertData(event, flatRaw, res)
	}
}

func HandleMail(event dto.CloudEvents, flatRaw dto.FlatRaw, pushConfig bo.RecipientPushConfig,
	push config.PushConfig) {
	if pushConfig.NeedMail {
		sendMailLog := "send mail %s %s %s"
		res := sendMail(event, pushConfig, push.EmailConfig)
		if res.Res == dto.Failed {
			logrus.Infof(sendMailLog, event.ID(), "failed", pushConfig.Mail)
		} else {
			logrus.Infof(sendMailLog, event.ID(), "success", pushConfig.Mail)
		}
		//insertData(event, flatRaw, res)
	}
}

func sendHWCloudMessage(raw dto.Raw, recipient bo.RecipientPushConfig,
	messageConfig pushSdk.MsgConfig) dto.PushResult {
	templateParas := raw.ToMessageArgs(recipient.MessageTemplate)
	res := pushSdk.SendHWCloudMessage(messageConfig, templateParas, recipient)
	if res.Res == dto.Failed {
		logrus.Error("send hwcloud message failed", templateParas)
	}
	return res
}

func sendInnerMessage(event dto.CloudEvents, recipient bo.RecipientPushConfig) dto.PushResult {
	return event.SendInnerMessage(recipient)
}

func sendTodoMessage(event dto.CloudEvents, recipient bo.RecipientPushConfig) dto.PushResult {
	return event.SendTodoMessage(recipient)
}

func sendFollowMessage(event dto.CloudEvents, recipient bo.RecipientPushConfig) dto.PushResult {
	return event.SendFollowMessage(recipient)
}

func sendMail(event dto.CloudEvents, recipient bo.RecipientPushConfig, emailConfig pushSdk.EmailConfig) dto.PushResult {
	return pushSdk.SendEmail(event.Extensions()["mailtitle"].(string),
		event.Extensions()["mailsummary"].(string), recipient, emailConfig)
}

func checkNil(stringifyMap map[string]string,
	eurBuildEvent dto.CloudEvents) {
	// 创建一个用于存储 nil 值的切片
	var nilFields []string

	if stringifyMap == nil {
		nilFields = append(nilFields, "Event Data")
	}
	if eurBuildEvent.Extensions()["sourceurl"] == nil {
		nilFields = append(nilFields, "Source URL")
	}
	if eurBuildEvent.Extensions()["user"] == nil {
		nilFields = append(nilFields, "User")
	}
	if eurBuildEvent.Extensions()["title"] == nil {
		nilFields = append(nilFields, "Title")
	}
	if eurBuildEvent.Extensions()["summary"] == nil {
		nilFields = append(nilFields, "Summary")
	}

	// 打印所有值为 nil 的字段
	if len(nilFields) > 0 {
		logrus.Infof("Nil fields: %v", nilFields)
	} else {
		logrus.Infof("No nil fields found.")
	}
}

func insertData(eurBuildEvent dto.CloudEvents, flatRaw dto.FlatRaw, result dto.PushResult) {
	stringifyMap := flatRaw.StringifyMap()
	insert := `insert into message_push_record (recipient_id, time_uuid, created_at, event_data, event_data_content_type,
                                 event_data_schema, event_id, event_source, event_source_url, event_spec_version,
                                 event_time, event_type, event_user, push_address, push_state, push_time, push_type,
                                 remark, title, summary)
values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
`
	checkNil(stringifyMap, eurBuildEvent)
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
		logrus.Errorf("insert data failed, err:%v", err)
		panic(nil)
		return
	}
}
