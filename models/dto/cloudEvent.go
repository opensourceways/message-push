package dto

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(event)
	var subscribePushConfigs []bo.RecipientConfig
	postgresql.DB().Raw(
		`select sc.mode_filter,
       rc.recipient_id,
       rc.mail,
       rc.message,
       rc.phone,
       pc.id,
       pc.need_message,
       pc.need_phone,
       pc.need_mail,
       pc.need_inner_message
from message_center.subscribe_config sc
         join message_center.push_config pc
              on sc.id = pc.subscribe_id
         join message_center.recipient_config rc on pc.recipient_id = rc.recipient_id
where sc.source = ?
  and sc.is_deleted = false
  and pc.is_deleted = false
  and rc.is_deleted = false`,
		event.Source(),
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}

func (event CloudEvents) SendInnerMessage(recipient bo.RecipientConfig) PushResult {
	innerMessageDO := do.InnerMessageDO{
		EventId:     event.ID(),
		Source:      event.Source(),
		RecipientId: recipient.RecipientId,
	}
	if postgresql.DB().Model(&innerMessageDO).Where("recipient_id=?", recipient.RecipientId, "source=?", innerMessageDO.Source, "event_id = ?", innerMessageDO.EventId).Updates(&innerMessageDO).RowsAffected == 0 {
		postgresql.DB().Create(&innerMessageDO)
	}

	return PushResult{
		Res:         Succeed,
		Time:        time.Now(),
		Remark:      "succeed",
		RecipientId: recipient.RecipientId,
		PushType:    "inner message",
		PushAddress: "",
	}
}
