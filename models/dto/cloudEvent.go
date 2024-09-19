package dto

import (
	"encoding/json"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/do"
	"github.com/todocoder/go-stream/stream"
)

const related_sql = `
	select distinct id       recipient_id,
                mail,
                message,
                phone,
                null::jsonb as mode_filter,
                false as need_message,
                false as need_phone,
                false as need_mail,
                false  as need_inner_message,
                null  as message_template,
                null  as mail_template
from message_center.recipient_config
where is_deleted is false
  and (recipient_name in ?
    or mail in ?
    or phone in ?
      or gitee_user_name in ?
      or user_id in ?
    )
`
const subscribe_sql = `
select distinct rc.id recipient_id,

                rc.mail,
                rc.message,
                rc.phone,
                sc.mode_filter,

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
    	and (recipient_name in ?
                      or mail in ?
                      or phone in ?
                      or gitee_user_name in ?
                      or user_id in ?
                     )
         left join message_center.push_template pt on sc.source = pt.source and sc.event_type = pt.event_type
where sc.source = ?
  and sc.event_type = ?
  and sc.is_deleted = false
  and pc.is_deleted = false
  and rc.is_deleted = false
`

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

func (event CloudEvents) GetRecipient() []bo.RecipientPushConfig {
	subscribePushConfigs := event.getSubscribeFromDB()
	relatedPushConfigs := event.getRelatedFromDB()
	return mergeRecipient(subscribePushConfigs, relatedPushConfigs)
}

func mergeRecipient(subscribe []bo.RecipientPushConfig, related []bo.RecipientPushConfig) []bo.RecipientPushConfig {
	return stream.Of(subscribe...).Concat(stream.Of(related...)).Distinct(func(item bo.RecipientPushConfig) any {
		return item.RecipientId
	}).ToSlice()
}

func (event CloudEvents) getRelatedFromDB() []bo.RecipientPushConfig {
	relatedUsers := strings.Split(event.Extensions()["relatedusers"].(string), ",")
	var subscribePushConfigs []bo.RecipientPushConfig
	postgresql.DB().Raw(
		related_sql,
		relatedUsers, relatedUsers, relatedUsers, relatedUsers, relatedUsers,
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}

func (event CloudEvents) getSubscribeFromDB() []bo.RecipientPushConfig {
	relatedUsers := strings.Split(event.Extensions()["relatedusers"].(string), ",")
	var subscribePushConfigs []bo.RecipientPushConfig
	postgresql.DB().Raw(
		subscribe_sql,
		relatedUsers, relatedUsers, relatedUsers, relatedUsers, relatedUsers, event.Source(), event.Type(),
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
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

func (event CloudEvents) SendInnerMessage(recipient bo.RecipientPushConfig) PushResult {
	innerMessageDO := do.InnerMessageDO{
		EventId:     event.ID(),
		Source:      event.Source(),
		RecipientId: recipient.RecipientId,
		IsRead:      false,
		IsSpecial:   false,
	}
	return SaveDb(innerMessageDO)
}
