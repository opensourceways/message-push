package dto

import (
	"encoding/json"
	flattener "github.com/anshal21/json-flattener"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/todocoder/go-stream/stream"
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
	event := cloudevents.NewEvent()
	event.SetID(raw.ID)
	event.SetSource(
		"https://eur.openeuler.openatom.cn/coprs/" + raw.Body.Owner + "/" + raw.Body.Pkg + "/build/" + strconv.Itoa(raw.Body.Build),
	)
	event.SetType("state:change")
	event.SetTime(raw.Headers.SentAt)
	event.SetDataContentType("application/json")
	event.SetDataSchema("eur:build_task")
	event.SetSpecVersion("1.0")
	err := event.SetData(cloudevents.ApplicationJSON, raw)
	if err != nil {
		return EurBuildEvent{}
	}
	return EurBuildEvent{event}
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
	return stream.Of(subscribePushConfigs...).Filter(
		func(item bo.SubscribePushConfig) bool {
			return eurBuildRaw.ModeFilter(item.ModeFilter)
		},
	).ToSlice()
}

func getSubscribeFromDB(event EurBuildEvent) []bo.SubscribePushConfig {
	var subscribePushConfigs []bo.SubscribePushConfig
	postgresql.DB().Table("message_center.cloud_event_message").Raw(
		`select sc.data_schema,
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
  and sc.data_schema = ?
  and sc.spec_version  = ?
group by sc.data_schema,
         sc.event_type,
         sc.spec_version,
         sc.mode_filter`,
		event.Type(), event.DataSchema(), event.SpecVersion(),
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}
