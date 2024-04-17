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
	transferErr := publishEurEvent(raw)
	if transferErr != nil {
		return transferErr
	}
	save(raw)
	return nil
}

func publishEurEvent(raw dto.EurBuildRaw) error {
	eurBuildEvent := raw.ToCloudEvent()
	sendErr := kafka.SendMsg("eur_build_event", &eurBuildEvent)
	if sendErr != nil {
		return sendErr
	}
	return nil
}

func save(raw dto.EurBuildRaw) {
	do := raw.ToCloudEventDO()
	if postgresql.DB().Model(&do).Where("event_id = ?", do.EventId).Updates(&do).RowsAffected == 0 {
		postgresql.DB().Create(&do)
	}
}
