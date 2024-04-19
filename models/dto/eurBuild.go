package dto

import (
	"encoding/json"
	flattener "github.com/anshal21/json-flattener"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"message-push/common/postgresql"
	"message-push/models/bo"
	"message-push/models/do"
	"strconv"
	"time"
)

type EurBuildRaw struct {
	Body struct {
		Build   int    `json:"build"`
		Chroot  string `json:"chroot"`
		Copr    string `json:"copr"`
		IP      string `json:"ip"`
		Owner   string `json:"owner"`
		Pid     int    `json:"pid"`
		Pkg     string `json:"pkg"`
		Status  int    `json:"status"`
		User    string `json:"user"`
		Version string `json:"version"`
		What    string `json:"what"`
		Who     string `json:"who"`
	} `json:"body"`
	Headers struct {
		FedoraMessagingSchema     string    `json:"fedora_messaging_schema"`
		FedoraMessagingSeverity   int       `json:"fedora_messaging_severity"`
		FedoraMessagingUserPackit bool      `json:"fedora_messaging_user_packit"`
		Priority                  int       `json:"priority"`
		SentAt                    time.Time `json:"sent-at"`
	} `json:"headers"`
	ID    string      `json:"id"`
	Queue interface{} `json:"queue"`
	Topic string      `json:"topic"`
}

func (raw *EurBuildRaw) ToCloudEvent() EurBuildEvent {
	newEvent := cloudevents.NewEvent()
	newEvent.SetSource("https://eur.openeuler.openatom.cn")
	newEvent.SetDataSchema(
		"https://eur.openeuler.openatom.cn/coprs/" + raw.Body.Owner + "/" + raw.Body.Pkg + "/build/" + strconv.Itoa(raw.Body.Build),
	)
	newEvent.SetType(raw.Topic)
	newEvent.SetTime(raw.Headers.SentAt)
	newEvent.SetDataContentType(cloudevents.ApplicationCloudEventsJSON)
	newEvent.SetSpecVersion(cloudevents.VersionV1)
	_ = newEvent.SetData(cloudevents.ApplicationJSON, raw)
	newEvent.SetID(raw.Topic + ":" + raw.ID)

	return EurBuildEvent{newEvent}
}

type EurBuildEvent struct {
	cloudevents.Event
}

func (event EurBuildEvent) ToCloudEventDO() do.MessageCloudEventDO {
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

func (event EurBuildEvent) Message() ([]byte, error) {
	return json.Marshal(event)
}

func (raw *EurBuildRaw) Message() ([]byte, error) {
	return json.Marshal(raw)
}

func (raw *EurBuildRaw) Flatten() map[string]interface{} {
	s, _ := json.Marshal(raw)
	flatJSON, _ := flattener.FlattenJSON(string(s), flattener.DotSeparator)
	flatMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flatJSON), &flatMap)
	return flatMap
}

func (raw *EurBuildRaw) ModeFilter(modeFilterJson datatypes.JSON) bool {
	s, _ := json.Marshal(raw)
	flatJSON, _ := flattener.FlattenJSON(string(s), flattener.DotSeparator)

	flatMap := make(map[string]string)
	_ = json.Unmarshal([]byte(flatJSON), &flatMap)
	modeFilterMap := make(map[string]string)
	_ = json.Unmarshal(modeFilterJson, &modeFilterMap)
	validate := validator.New()
	for k, v := range modeFilterMap {
		err := validate.Var(flatMap[k], v)
		if err != nil {
			return false
		}
	}
	return true
}

func (event EurBuildEvent) GetSubscribe() []bo.SubscribePushConfig {
	subscribePushConfigs := getSubscribeFromDB(event)
	var eurBuildRaw EurBuildRaw
	_ = json.Unmarshal(event.Data(), &eurBuildRaw)
	return subscribePushConfigs
}

func getSubscribeFromDB(event EurBuildEvent) []bo.SubscribePushConfig {
	var subscribePushConfigs []bo.SubscribePushConfig
	postgresql.DB().Table("message_center.cloud_event_message").Raw(
		`select sc.source,
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
group by sc.source,
         sc.event_type,
         sc.spec_version,
         sc.mode_filter`,
		event.Type(), event.Source(), event.SpecVersion(),
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}
