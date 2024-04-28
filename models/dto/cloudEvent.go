package dto

import (
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"message-push/common/postgresql"
	"message-push/models/bo"
	"message-push/models/do"
)

type CloudEvents struct {
	cloudevents.Event
}

func (event CloudEvents) Message() ([]byte, error) {
	return json.Marshal(event)
}

func (event CloudEvents) ToCloudEventDO() do.MessageCloudEventDO {
	messageCloudEventDO := do.MessageCloudEventDO{
		Source:          event.Source(),
		Time:            event.Time(),
		EventType:       event.Type(),
		SpecVersion:     event.SpecVersion(),
		DataSchema:      event.DataSchema(),
		DataContentType: event.DataContentType(),
		EventId:         event.ID(),
		DataJson:        event.Data(),
	}
	return messageCloudEventDO
}

func (event CloudEvents) GetSubscribe() []bo.SubscribePushConfig {
	subscribePushConfigs := event.getSubscribeFromDB()
	return subscribePushConfigs
}

func (event CloudEvents) getSubscribeFromDB() []bo.SubscribePushConfig {
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

func (event CloudEvents) SaveDb() {
	do := event.ToCloudEventDO()
	if postgresql.DB().Model(&do).Where("source=?", do.Source, "event_id = ?", do.EventId).Updates(&do).RowsAffected == 0 {
		postgresql.DB().Create(&do)
	}
}
