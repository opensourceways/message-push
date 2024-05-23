package dto

import (
	"encoding/json"
	flattener "github.com/anshal21/json-flattener"
)

type EurBuildMessageRaw struct {
	Body struct {
		User    string      `json:"user"`
		Copr    string      `json:"copr"`
		Owner   string      `json:"owner"`
		Pkg     interface{} `json:"pkg"`
		Build   int         `json:"build"`
		Chroot  string      `json:"chroot"`
		Version interface{} `json:"version"`
		Status  int         `json:"status"`
		IP      string      `json:"ip"`
		Who     string      `json:"who"`
		Pid     int         `json:"pid"`
		What    string      `json:"what"`
	} `json:"body"`
	Headers struct {
		OpenEulerMessagingSchema string `json:"openEuler_messaging_schema"`
		SentAt                   string `json:"sent-at"`
	} `json:"headers"`
	ID    string `json:"id"`
	Topic string `json:"topic"`
}

func (raw *EurBuildMessageRaw) Flatten() map[string]interface{} {
	s, _ := json.Marshal(raw)
	flatJSON, _ := flattener.FlattenJSON(string(s), flattener.DotSeparator)
	flatMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flatJSON), &flatMap)
	return flatMap
}
