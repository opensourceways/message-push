package dto

import (
	"encoding/json"
	flattener "github.com/anshal21/json-flattener"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"message-push/common/postgresql"
	"message-push/models/bo"
	"message-push/models/do"
	"strconv"
	"time"
)

type EurBuildMessageRaw struct {
	Properties struct {
		AppID           interface{} `json:"app_id"`
		ClusterID       interface{} `json:"cluster_id"`
		ContentEncoding string      `json:"content_encoding"`
		ContentType     string      `json:"content_type"`
		CorrelationID   interface{} `json:"correlation_id"`
		DeliveryMode    int         `json:"delivery_mode"`
		Expiration      interface{} `json:"expiration"`
		Headers         struct {
			FedoraMessagingSchema      string    `json:"fedora_messaging_schema"`
			FedoraMessagingSeverity    int       `json:"fedora_messaging_severity"`
			FedoraMessagingUserKkleine bool      `json:"fedora_messaging_user_kkleine"`
			Priority                   int       `json:"priority"`
			SentAt                     time.Time `json:"sent-at"`
			XReceivedFrom              []struct {
				ClusterName string `json:"cluster-name"`
				Exchange    string `json:"exchange"`
				Redelivered bool   `json:"redelivered"`
				URI         string `json:"uri"`
			} `json:"x-received-from"`
		} `json:"headers"`
		MessageID string      `json:"message_id"`
		Priority  interface{} `json:"priority"`
		ReplyTo   interface{} `json:"reply_to"`
		Timestamp interface{} `json:"timestamp"`
		Type      interface{} `json:"type"`
		UserID    interface{} `json:"user_id"`
	} `json:"_properties"`
	Body struct {
		Build   int    `json:"build"`
		Chroot  string `json:"chroot"`
		Copr    string `json:"copr"`
		IP      string `json:"ip"`
		Owner   string `json:"owner"`
		PID     int    `json:"pid"`
		Pkg     string `json:"pkg"`
		Status  int    `json:"status"`
		User    string `json:"user"`
		Version string `json:"version"`
		What    string `json:"what"`
		Who     string `json:"who"`
	} `json:"body"`
	Queue    string `json:"queue"`
	Severity int    `json:"severity"`
	Topic    string `json:"topic"`
}

func (raw *EurBuildMessageRaw) Flatten() map[string]interface{} {
	s, _ := json.Marshal(raw)
	flatJSON, _ := flattener.FlattenJSON(string(s), flattener.DotSeparator)
	flatMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flatJSON), &flatMap)
	return flatMap
}

func (raw *EurBuildMessageRaw) ToCloudEvent() EurBuildEvent {
	newEvent := cloudevents.NewEvent()
	newEvent.SetSource("https://copr.fedorainfracloud.org")
	newEvent.SetDataSchema(
		"https://copr.fedorainfracloud.org/" + raw.Body.Owner + "/" + raw.Body.Pkg + "/build/" + strconv.Itoa(raw.Body.Build),
	)
	newEvent.SetType(raw.Topic)
	newEvent.SetTime(raw.Properties.Headers.SentAt)
	newEvent.SetDataContentType(cloudevents.ApplicationCloudEventsJSON)
	newEvent.SetSpecVersion(cloudevents.VersionV1)
	_ = newEvent.SetData(cloudevents.ApplicationJSON, raw)
	newEvent.SetID(raw.Topic + ":" + raw.Properties.MessageID)

	return EurBuildEvent{newEvent}
}

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

func (event EurBuildEvent) GetSubscribe() []bo.SubscribePushConfig {
	subscribePushConfigs := getSubscribeFromDB(event)
	return subscribePushConfigs
}

func getSubscribeFromDB(event EurBuildEvent) []bo.SubscribePushConfig {
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
