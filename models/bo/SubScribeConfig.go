package bo

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type SubscribePushConfig struct {
	DataSchema  string          `gorm:"column:data_schema"`
	EventType   string          `gorm:"column:event_type"`
	SpecVersion string          `gorm:"column:spec_version"`
	PushConfigs json.RawMessage `gorm:"type:json;column:push_configs" json:"push_configs"`
	ModeFilter  datatypes.JSON  `gorm:"column:mode_filter"`
}

type PushConfig struct {
	PushType    string `gorm:"column:pushType" json:"push_type"`
	PushAddress string `gorm:"column:push_address" json:"push_address"`
}
