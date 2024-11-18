package do

import (
	"github.com/opensourceways/message-push/common/postgresql"
)

type InnerMessageDO struct {
	postgresql.CommonModel
	EventId     string `gorm:"column:event_id"`
	Source      string `gorm:"column:source"`
	RecipientId string `gorm:"column:recipient_id"`
	IsRead      bool   `gorm:"column:is_read"`
	IsSpecial   bool   `gorm:"column:is_special"`
}

type TodoMessageDO struct {
	postgresql.CommonModel
	EventId     string `gorm:"column:event_id"`
	Source      string `gorm:"column:source"`
	RecipientId string `gorm:"column:recipient_id"`
	IsRead      bool   `gorm:"column:is_read"`
	IsSpecial   bool   `gorm:"column:is_special"`
	IsDone      bool   `gorm:"column:is_done"`
}

func (m *InnerMessageDO) TableName() string {
	return "message_center.inner_message"
}
