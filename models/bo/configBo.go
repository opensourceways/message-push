package bo

import (
	"gorm.io/datatypes"
)

type RecipientPushConfig struct {
	RecipientId      string         `gorm:"column:recipient_id"`
	Mail             string         `gorm:"column:mail"`
	Message          string         `gorm:"column:message"`
	Phone            string         `gorm:"column:phone"`
	ModeFilter       datatypes.JSON `gorm:"column:mode_filter"`
	NeedMessage      bool           `gorm:"column:need_message"`
	NeedPhone        bool           `gorm:"column:need_phone"`
	NeedMail         bool           `gorm:"column:need_mail"`
	NeedInnerMessage bool           `gorm:"column:need_inner_message"`
	MessageTemplate  string         `gorm:"column:message_template"`
	MailTemplate     string         `gorm:"column:mail_template"`
	IsSpecial        bool           `gorm:"column:is_special"`
}
