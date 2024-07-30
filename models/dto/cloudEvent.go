package dto

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/do"
	"time"
)

type CloudEvents struct {
	cloudevents.Event
}

func NewCloudEvents() CloudEvents {
	return CloudEvents{
		Event: cloudevents.NewEvent(cloudevents.VersionV1),
	}
}

func (event CloudEvents) Message() ([]byte, error) {
	return json.Marshal(event)
}

func (event CloudEvents) GetRecipient() []bo.RecipientConfig {
	subscribePushConfigs := event.getRecipientFromDB()
	return subscribePushConfigs
}

func (event CloudEvents) getRecipientFromDB() []bo.RecipientConfig {
	var subscribePushConfigs []bo.RecipientConfig
	postgresql.DB().Raw(
		`select  distinct sc.mode_filter,
       rc.id recipient_id,
       rc.mail,
       rc.message,
       rc.phone,
       pc.id,
       pc.need_message,
       pc.need_phone,
       pc.need_mail,
       pc.need_inner_message,
       pt.message_template,
       pt.mail_template
from message_center.subscribe_config sc
         join message_center.push_config pc
              on sc.id = pc.subscribe_id
         join message_center.recipient_config rc on pc.recipient_id = rc.id
         left join message_center.push_template pt on sc.source = pt.source and sc.event_type = pt.event_type
where sc.source = ?
  and sc.event_type = ?
  and sc.is_deleted = false
  and pc.is_deleted = false
  and rc.is_deleted = false`,
		event.Source(), event.Type(),
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}

func (event CloudEvents) SendInnerMessageByRelatedUsers(relatedUsers []string) {
	for _, user := range relatedUsers {
		innerMessageDO := do.InnerMessageDO{
			EventId:     event.ID(),
			Source:      event.Source(),
			RecipientId: user,
			IsRead:      false,
		}
		SaveDb(innerMessageDO)
	}
}

func SaveDb(m do.InnerMessageDO) PushResult {
	res := postgresql.DB().Save(&m)
	if res.Error != nil {
		return PushResult{
			Res:         Failed,
			Time:        time.Now(),
			Remark:      res.Error.Error(),
			RecipientId: m.RecipientId,
			PushType:    "inner message",
			PushAddress: "",
		}
	} else {
		return PushResult{
			Res:         Succeed,
			Time:        time.Now(),
			Remark:      "succeed",
			RecipientId: m.RecipientId,
			PushType:    "inner message",
			PushAddress: "",
		}
	}
}

func (event CloudEvents) SendInnerMessage(recipient bo.RecipientConfig) PushResult {
	innerMessageDO := do.InnerMessageDO{
		EventId:     event.ID(),
		Source:      event.Source(),
		RecipientId: recipient.RecipientId,
		IsRead:      false,
	}
	return SaveDb(innerMessageDO)
}
