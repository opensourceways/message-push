package do

import (
	"gorm.io/datatypes"
	"message-push/common/postgresql"
	"time"
)

type JSONB []interface{}

type MessageCloudEventDO struct {
	postgresql.CommonModel
	EventId         string         `gorm:"column:event_id"`
	Source          string         `gorm:"column:source"`
	EventType       string         `gorm:"column:type"`
	DataContentType string         `gorm:"column:data_content_type"`
	DataSchema      string         `gorm:"column:data_schema"`
	SpecVersion     string         `gorm:"column:spec_version"`
	Time            time.Time      `gorm:"column:time"`
	DataJson        datatypes.JSON `gorm:"column:data_json"`
}

func (m *MessageCloudEventDO) TableName() string {
	return "cloud_event_message"
}

type SubScribeConfigDO struct {
	postgresql.CommonModel
	Source     string         `gorm:"column:source"`
	EventType  string         `gorm:"column:type"`
	Version    string         `gorm:"column:version"`
	ModeFilter datatypes.JSON `gorm:"column:mod_filter"`
	IsDeleted  bool           `gorm:"column:is_deleted"`
}

func (m *SubScribeConfigDO) TableName() string {
	return "subscribe_config"
}

type PushConfigDO struct {
	postgresql.CommonModel
	SubScribeId int    `gorm:"column:subscribe_id"`
	PushType    string `gorm:"column:type"`
	PushAddress string `gorm:"column:version"`
	IsDeleted   bool   `gorm:"column:is_deleted"`
}

func (m *SubScribeConfigDO) PushConfigDO() string {
	return "push_config"
}
