package service

import (
	"github.com/goccy/go-json"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
)

func GiteeIssueHandle(payload []byte, _ map[string]string) error {
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

func handle(event dto.CloudEvents, push config.PushConfig) error {
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
			return utils.ModeFilter(flatRaw, recipient.ModeFilter)
		},
	).ForEach(func(recipient bo.RecipientConfig) {
		if recipient.NeedMessage {
			res := sendHWCloudMessage(raw, recipient, push.MsgConfig)
			pushSdk.InsertData(event, flatRaw, res)
			logrus.Info("send message ", event.ID()+" success")
		}
		if recipient.NeedMail {
			res := sendMail(event, recipient, push.EmailConfig)
			pushSdk.InsertData(event, flatRaw, res)
			logrus.Info("send mail ", event.ID()+" success")
		}
		if recipient.NeedInnerMessage {
			res := sendInnerMessage(event, recipient)
			logrus.Info("send inner message ", event.ID()+" success")
			pushSdk.InsertData(event, flatRaw, res)
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
