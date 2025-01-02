package dto

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/sirupsen/logrus"
	"github.com/todocoder/go-stream/stream"
	"gorm.io/gorm/clause"

	"github.com/opensourceways/message-push/common/postgresql"
	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/do"
)

const related_sql = `
	select distinct id       recipient_id,
                mail,
                message,
                phone,
                null::jsonb as mode_filter,
                false as need_message,
                false as need_phone,
                false as need_mail,
                false  as need_inner_message,
                null  as message_template,
                null  as mail_template
from message_center.recipient_config
where is_deleted is false
  and (recipient_name in ?
    or mail in ?
    or phone in ?
      or gitee_user_name in ?
      or user_id in ?
    )
`
const subscribe_sql = `
select distinct rc.id recipient_id,

                rc.mail,
                rc.message,
                rc.phone,
                sc.mode_filter,

                pc.need_message,
                pc.need_phone,
                pc.need_mail,
                pc.need_inner_message,
                pt.message_template,
                pt.mail_template
from message_center.subscribe_config sc
         join message_center.push_config pc
              on sc.id = pc.subscribe_id
         join message_center.recipient_config rc on pc.recipient_id = rc.id
    	and (recipient_name in ?
                      or mail in ?
                      or phone in ?
                      or gitee_user_name in ?
                      or user_id in ?
                     )
         left join message_center.push_template pt on sc.source = pt.source and sc.event_type = pt.event_type
where sc.source = ?
  and sc.event_type = ?
  and sc.is_deleted = false
  and pc.is_deleted = false
  and rc.is_deleted = false
`

type CloudEvents struct {
	cloudevents.Event
}

func NewCloudEvents() CloudEvents {
	return CloudEvents{
		Event: cloudevents.NewEvent(cloudevents.VersionV1),
	}
}

func (event CloudEvents) Message() ([]byte, error) {
	return json.Marshal(event)
}

func (event CloudEvents) GetRecipient() []bo.RecipientPushConfig {
	relatedPushConfigs := event.GetRelatedFromDB()
	subscribePushConfigs := event.GetSubscribeFromDB()
	return mergeRecipient(subscribePushConfigs, relatedPushConfigs)
}

func mergeRecipient(subscribe []bo.RecipientPushConfig, related []bo.RecipientPushConfig) []bo.RecipientPushConfig {
	var unique []string
	subs := stream.Of(subscribe...).Distinct(func(item bo.
		RecipientPushConfig) any {
		return fmt.Sprintf("%s:%v", item.RecipientId, item.ModeFilter)
	}).ToSlice()
	for _, sub := range subs {
		if !slices.Contains(unique, sub.RecipientId) {
			logrus.Infof("the mail is %v, subs config is %v", sub.Mail, sub.NeedMail)
			unique = append(unique, sub.RecipientId)
		}
	}
	for _, relate := range related {
		if slices.Contains(unique, relate.RecipientId) {
			continue
		}
		unique = append(unique, relate.RecipientId)
		subs = append(subs, relate)
	}
	return subs
}

func (event CloudEvents) GetRelatedFromDB() []bo.RecipientPushConfig {
	relatedUsers, ok := event.Extensions()["relatedusers"].(string)
	logrus.Infof("relatedUsers :%v", relatedUsers)
	if !ok || relatedUsers == "" {
		return nil
	}
	relatedUsersList := strings.Split(relatedUsers, ",")
	var subscribePushConfigs []bo.RecipientPushConfig
	postgresql.DB().Raw(
		related_sql,
		relatedUsersList, relatedUsersList, relatedUsersList, relatedUsersList, relatedUsersList,
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}

func (event CloudEvents) GetTodoFromDB() []bo.RecipientPushConfig {
	todoUsers, ok := event.Extensions()["todousers"].(string)
	if !ok || todoUsers == "" {
		return nil
	}
	todoUsersList := strings.Split(todoUsers, ",")
	var todoPushConfigs []bo.RecipientPushConfig
	postgresql.DB().Raw(
		related_sql,
		todoUsersList, todoUsersList, todoUsersList, todoUsersList, todoUsersList,
	).Scan(&todoPushConfigs)
	return todoPushConfigs
}

func (event CloudEvents) GetFollowFromDB() []bo.RecipientPushConfig {
	followUsers, ok := event.Extensions()["followusers"].(string)
	if !ok || followUsers == "" {
		return nil
	}
	followUsersList := strings.Split(followUsers, ",")
	var followPushConfigs []bo.RecipientPushConfig
	postgresql.DB().Raw(
		related_sql,
		followUsersList, followUsersList, followUsersList, followUsersList, followUsersList,
	).Scan(&followPushConfigs)
	return followPushConfigs
}

func (event CloudEvents) GetSubscribeFromDB() []bo.RecipientPushConfig {
	if event.Extensions()["relatedusers"] == nil {
		return nil
	}
	relatedUsers := strings.Split(event.Extensions()["relatedusers"].(string), ",")
	var subscribePushConfigs []bo.RecipientPushConfig
	postgresql.DB().Raw(
		subscribe_sql,
		relatedUsers, relatedUsers, relatedUsers, relatedUsers, relatedUsers, event.Source(), event.Type(),
	).Scan(&subscribePushConfigs)
	return subscribePushConfigs
}

func SaveDb(m do.RelatedMessageDO) PushResult {
	res := postgresql.DB().Save(&m)
	if res.Error != nil {
		return PushResult{
			Res:         Failed,
			Time:        time.Now(),
			Remark:      res.Error.Error(),
			RecipientId: m.RecipientId,
			PushType:    "related message",
			PushAddress: "",
		}
	} else {
		return PushResult{
			Res:         Succeed,
			Time:        time.Now(),
			Remark:      "succeed",
			RecipientId: m.RecipientId,
			PushType:    "related message",
			PushAddress: "",
		}
	}
}

func (event CloudEvents) SendInnerMessage(recipient bo.RecipientPushConfig) PushResult {
	innerMessageDO := do.RelatedMessageDO{
		EventId:     event.ID(),
		Source:      event.Source(),
		RecipientId: recipient.RecipientId,
		IsRead:      false,
	}
	return SaveDb(innerMessageDO)
}

func SaveTodoDb(m do.TodoMessageDO) PushResult {
	res := postgresql.DB().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "business_id"}, {Name: "recipient_id"}, {Name: "source"}},
		DoUpdates: clause.AssignmentColumns([]string{"latest_event_id", "is_done", "updated_at"}),
	}).Create(&m)
	if res.Error != nil {
		return PushResult{
			Res:         Failed,
			Time:        time.Now(),
			Remark:      res.Error.Error(),
			RecipientId: m.RecipientId,
			PushType:    "todo message",
			PushAddress: "",
		}
	} else {
		return PushResult{
			Res:         Succeed,
			Time:        time.Now(),
			Remark:      "succeed",
			RecipientId: m.RecipientId,
			PushType:    "todo message",
			PushAddress: "",
		}
	}
}

func (event CloudEvents) SendTodoMessage(recipient bo.RecipientPushConfig) PushResult {
	businessid, ok := event.Extensions()["businessid"].(string)
	if !ok || businessid == "" {
		businessid = ""
	}
	isDone, ok := event.Extensions()["isdone"].(bool)
	if !ok {
		isDone = false
	}

	todoMessageDO := do.TodoMessageDO{
		BusinessId:    businessid,
		Source:        event.Source(),
		RecipientId:   recipient.RecipientId,
		LatestEventId: event.ID(),
		IsDone:        isDone,
	}
	return SaveTodoDb(todoMessageDO)
}

func SaveFollowDb(m do.FollowMessageDO) PushResult {
	res := postgresql.DB().Save(&m)
	if res.Error != nil {
		return PushResult{
			Res:         Failed,
			Time:        time.Now(),
			Remark:      res.Error.Error(),
			RecipientId: m.RecipientId,
			PushType:    "follow message",
			PushAddress: "",
		}
	} else {
		return PushResult{
			Res:         Succeed,
			Time:        time.Now(),
			Remark:      "succeed",
			RecipientId: m.RecipientId,
			PushType:    "follow message",
			PushAddress: "",
		}
	}
}

func (event CloudEvents) SendFollowMessage(recipient bo.RecipientPushConfig) PushResult {
	followMessageDO := do.FollowMessageDO{
		EventId:     event.ID(),
		Source:      event.Source(),
		RecipientId: recipient.RecipientId,
		IsRead:      false,
	}
	return SaveFollowDb(followMessageDO)
}
