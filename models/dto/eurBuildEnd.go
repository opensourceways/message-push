package dto

import (
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
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

func NewEurBuildRaw() EurBuildRaw {
	EurBuildJSON := `{
  "body": {
    "build": 7279434,
    "chroot": "fedora-39-x86_64",
    "copr": "cran",
    "ip": "2620:52:3:1:dead:beef:cafe:c156",
    "owner": "iucar",
    "pid": 1961158,
    "pkg": "R-CRAN-shortIRT",
    "status": 3,
    "user": "iucar",
    "version": "0.1.3-1.copr7279434",
    "what": "build start: user:iucar copr:cran pkg:R-CRAN-shortIRT build:7279434 ip:2620:52:3:1:dead:beef:cafe:c156 pid:1961158",
    "who": "backend.worker-rpm_build_worker:7279434-fedora-39-x86_64"
  },
  "headers": {
    "fedora_messaging_schema": "copr.build.start",
    "fedora_messaging_severity": 20,
    "fedora_messaging_user_iucar": true,
    "priority": 0,
    "sent-at": "2024-04-09T07:44:31+00:00"
  },
  "id": "d4b3c30c-c7f4-454a-ab0b-def09796bd90",
  "queue": null,
  "topic": "org.fedoraproject.prod.copr.build.start"
}`

	var raw EurBuildRaw

	err := json.Unmarshal([]byte(EurBuildJSON), &raw)
	if err != nil {
		return EurBuildRaw{}
	}

	return raw
}

type EurBuildEvent struct {
	cloudevents.Event
}

func (raw *EurBuildRaw) ToCloudEvent() EurBuildEvent {
	event := cloudevents.NewEvent()
	event.SetID(raw.ID)
	event.SetSource(
		"https://eur.openeuler.openatom.cn/coprs/" + raw.Body.Owner + "/" + raw.Body.Pkg + "/build/" + strconv.Itoa(raw.Body.Build),
	)
	event.SetType("state:change")
	event.SetTime(time.Now())
	event.SetDataContentType("application/json")
	event.SetDataSchema("eur:build_task")
	event.SetSpecVersion("0.0.1")
	err := event.SetData(cloudevents.ApplicationJSON, raw)
	if err != nil {
		return EurBuildEvent{}

	}
	return EurBuildEvent{event}
}

func (raw *EurBuildRaw) ToCloudEventDO() do.MessageCloudEventDO {
	jsons, errs := json.Marshal(raw) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	messageCloudEventDO := do.MessageCloudEventDO{
		Source:          "https://eur.openeuler.openatom.cn/coprs/" + raw.Body.Owner + "/" + raw.Body.Pkg + "/build/" + strconv.Itoa(raw.Body.Build),
		Time:            time.Now(),
		EventType:       "state:change",
		SpecVersion:     "0.0.1",
		DataSchema:      "eur:build_task",
		DataContentType: "application/json",
		EventId:         raw.ID,
		DataJson:        jsons,
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
  and sc.spec_version = ?
group by sc.data_schema,
         sc.event_type,
         sc.spec_version,
         sc.mode_filter`,
		event.Type(), event.DataSchema(), "0.0.1",
	).Scan(&subscribePushConfigs)

	return subscribePushConfigs
}
