package dto

import (
	"encoding/json"
	flattener "github.com/anshal21/json-flattener"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"message-push/common/postgresql"
	"message-push/models/bo"
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
	json.Unmarshal([]byte(flatJSON), &flatMap)

	modeFilterMap := make(map[string]string)

	json.Unmarshal(modeFilterJson, &modeFilterMap)
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
	return subscribePushConfigs
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
