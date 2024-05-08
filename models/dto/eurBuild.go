package dto

import (
	"encoding/json"
	flattener "github.com/anshal21/json-flattener"
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
