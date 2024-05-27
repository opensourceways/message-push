package bo

import (
	"gorm.io/datatypes"
)

type RecipientConfig struct {
	RecipientId      string         `gorm:"column:recipient_id"`
	Mail             string         `gorm:"column:mail"`
	Message          string         `gorm:"column:message"`
	Phone            string         `gorm:"column:phone"`
	ModeFilter       datatypes.JSON `gorm:"column:mode_filter"`
	NeedMessage      bool           `gorm:"column:need_message"`
	NeedPhone        bool           `gorm:"column:need_phone"`
	NeedMail         bool           `gorm:"column:need_mail"`
	NeedInnerMessage bool           `gorm:"column:need_inner_message"`
}
