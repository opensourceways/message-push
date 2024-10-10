package dto

import "time"

const Succeed = "succeed"
const Failed = "failed"

type PushResult struct {
	RecipientId string `json:"recipient_id"`
	PushAddress string `json:"push_address"`
	PushType    string `json:"push_type"`
	Res         string
	Remark      string
	Time        time.Time
}
