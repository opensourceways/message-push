package do

import (
	"github.com/opensourceways/message-push/common/postgresql"
)

type InnerMessageDO struct {
	postgresql.CommonModel
	EventId     string `gorm:"column:event_id"`
	Source      string `gorm:"column:source"`
	RecipientId string `gorm:"column:recipient_id"`
}

func (m *InnerMessageDO) TableName() string {
	return "test_message_center.message_center.inner_message"
}
