package dto

import (
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/models/bo"
	do "github.com/opensourceways/message-push/models/do"
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

func (event CloudEvents) GetSubscribe() []bo.SubscribePushConfig {
	subscribePushConfigs := event.getSubscribeFromDB()
	return subscribePushConfigs
}

func (event CloudEvents) getSubscribeFromDB() []bo.SubscribePushConfig {
	fmt.Println(event)
	var subscribePushConfigs []bo.SubscribePushConfig
	postgresql.DB().Raw(
		`select 
       sc.recipient_id,
       sc.source,
       sc.event_type,
       sc.spec_version,
       sc.mode_filter,
       json_agg(
               json_build_object('push_type', pc.push_type, 'push_address', pc.push_address)
       ) push_configs

from message_center.push_config pc
         join message_center.subscribe_config sc on pc.subscribe_id = sc.id
    and pc.is_deleted is false and sc.is_deleted = false
where sc.event_type = ?
  and sc.source = ?
  and sc.spec_version  = ?
  and sc.is_deleted = false
group by 
		sc.recipient_id,
		sc.source,
		sc.event_type,
		sc.spec_version,
		sc.mode_filter`,
		event.Type(), event.Source(), event.SpecVersion(),
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}

func (event CloudEvents) SendInnerMessage(recipientId string) PushResult {
	innerMessageDO := do.InnerMessageDO{
		EventId:     event.ID(),
		Source:      event.Source(),
		RecipientId: recipientId,
	}
	if postgresql.DB().Model(&innerMessageDO).Where("recipient_id=?", recipientId, "source=?", innerMessageDO.Source, "event_id = ?", innerMessageDO.EventId).Updates(&innerMessageDO).RowsAffected == 0 {
		postgresql.DB().Create(&innerMessageDO)
	}

	return PushResult{
		Res:    Succeed,
		Time:   time.Now(),
		Remark: "succeed",
	}
}
