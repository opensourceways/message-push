package do

import (
	"github.com/opensourceways/message-push/common/postgresql"
)

type RelatedMessageDO struct {
	postgresql.CommonModel
	EventId     string `gorm:"column:event_id"`
	Source      string `gorm:"column:source"`
	RecipientId string `gorm:"column:recipient_id"`
	IsRead      bool   `gorm:"column:is_read"`
}

func (m *RelatedMessageDO) TableName() string {
	return "message_center.related_message"
}

type TodoMessageDO struct {
	postgresql.CommonModel
	BusinessId    string `gorm:"column:business_id"`
	Source        string `gorm:"column:source"`
	RecipientId   string `gorm:"column:recipient_id"`
	LatestEventId string `gorm:"column:latest_event_id"`
	IsDone        bool   `gorm:"column:is_done"`
}

func (m *TodoMessageDO) TableName() string {
	return "message_center.todo_message"
}

type FollowMessageDO struct {
	postgresql.CommonModel
	EventId     string `gorm:"column:event_id"`
	Source      string `gorm:"column:source"`
	RecipientId string `gorm:"column:recipient_id"`
	IsRead      bool   `gorm:"column:is_read"`
}

func (m *FollowMessageDO) TableName() string {
	return "message_center.follow_message"
}
