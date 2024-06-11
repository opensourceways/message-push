package dto

import (
	"encoding/json"
	flattener "github.com/anshal21/json-flattener"
	"github.com/opensourceways/go-gitee/gitee"
)

type GiteeIssueRaw struct {
	gitee.IssueEvent
}

func (raw *GiteeIssueRaw) Flatten() map[string]interface{} {
	s, _ := json.Marshal(raw)
	flatJSON, _ := flattener.FlattenJSON(string(s), flattener.DotSeparator)
	flatMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flatJSON), &flatMap)
	return flatMap
}
