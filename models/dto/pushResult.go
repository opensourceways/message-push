package dto

import "time"

const Succeed = "succeed"
const Failed = "failed"

type PushResult struct {
	Res    string
	Remark string
	Time   time.Time
}
