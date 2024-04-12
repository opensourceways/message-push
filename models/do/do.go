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
