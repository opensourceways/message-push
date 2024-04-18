package transfer

import (
	"encoding/json"
	"message-push/common/kafka"
	"message-push/common/postgresql"
	"message-push/models/dto"
)

func Handle(payload []byte, _ map[string]string) error {
	var raw dto.EurBuildRaw
	msgBodyErr := json.Unmarshal(payload, &raw)
	if msgBodyErr != nil {
		return msgBodyErr
	}
	eurBuildEvent := raw.ToCloudEvent()
	kafkaSendErr := kafka.SendMsg("eur_build_event", &eurBuildEvent)
	if kafkaSendErr != nil {
		return kafkaSendErr
	}
	save(eurBuildEvent)
	return nil
}

func save(event dto.EurBuildEvent) {
	do := event.ToCloudEventDO()
	if postgresql.DB().Model(&do).Where("source=?", do.Source, "event_id = ?", do.EventId).Updates(&do).RowsAffected == 0 {
		postgresql.DB().Create(&do)
	}
}
