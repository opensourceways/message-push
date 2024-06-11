package push

import (
	"encoding/json"
	"github.com/opensourceways/message-push/common/pushSdk"
	"github.com/opensourceways/message-push/config"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
)

type GiteeIssueHandler struct{}

func GiteeIssueHandle(payload []byte, _ map[string]string) error {
	var handler GiteeIssueHandler

	giteeIssueEvent := dto.NewCloudEvents()
	msgBodyErr := json.Unmarshal(payload, &giteeIssueEvent)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	handler.publishMessage(giteeIssueEvent)
	return nil
}

func (handler *GiteeIssueHandler) publishMessage(event dto.CloudEvents) {
	var raw dto.GiteeIssueRaw
	_ = json.Unmarshal(event.Data(), &raw)
	recipients := event.GetRecipient()
	if recipients == nil || len(recipients) == 0 {
		return
	}
	flatRaw := raw.Flatten()
	stream.Of(recipients...).Filter(
		func(recipient bo.RecipientConfig) bool {
			if recipient.ModeFilter == nil {
				return true
			}
			return utils.ModeFilter(flatRaw, recipient.ModeFilter)
		},
	).ForEach(func(recipient bo.RecipientConfig) {
		if recipient.NeedMessage {
			res := handler.sendHWCloudMessage(raw, recipient)
			pushSdk.InsertData(event, flatRaw, res)
			logrus.Info("send message ", event.ID()+" success")
		}
		if recipient.NeedMail {
			res := handler.sendMail(event, recipient)
			pushSdk.InsertData(event, flatRaw, res)
			logrus.Info("send mail ", event.ID()+" success")
		}
		if recipient.NeedInnerMessage {
			res := handler.sendInnerMessage(event, recipient)
			logrus.Info("send inner message ", event.ID()+" success")
			pushSdk.InsertData(event, flatRaw, res)
		}
	})
}

func (handler *GiteeIssueHandler) sendHWCloudMessage(raw dto.GiteeIssueRaw, recipient bo.RecipientConfig) dto.PushResult {

	var templateParas []string
	return pushSdk.SendHWCloudMessage(config.GiteeConfigInstance.HWCloudMsgConfig, templateParas, recipient)
}

func (handler *GiteeIssueHandler) sendInnerMessage(event dto.CloudEvents, recipient bo.RecipientConfig) dto.PushResult {
	return event.SendInnerMessage(recipient)
}

func (handler *GiteeIssueHandler) sendMail(event dto.CloudEvents, recipient bo.RecipientConfig) dto.PushResult {
	return pushSdk.SendEmail(event.Extensions()["title"].(string),
		event.Extensions()["summary"].(string), recipient, config.EurBuildConfigInstance.EmailConfig)
}
