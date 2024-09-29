package pushSdk

import (
	"crypto/tls"
	"log"
	"time"

	"gopkg.in/mail.v2"

	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
)

type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPassword string `json:"smtp_password"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPSender   string `json:"smtp_sender"`
	SMTPUsername string `json:"smtp_username"`
}

func SendEmail(title string, summary string, recipient bo.RecipientPushConfig, config EmailConfig) dto.PushResult {
	err := sendSSLEmail(recipient.Mail, title, summary, config)
	if err != nil {
		return dto.PushResult{Res: dto.Failed, Remark: err.Error()}
	}
	return dto.PushResult{
		Res:         dto.Succeed,
		Time:        time.Now(),
		Remark:      "succeed",
		RecipientId: recipient.RecipientId,
		PushAddress: recipient.Mail,
		PushType:    "mail",
	}
}

func sendSSLEmail(receiver, subject, htmlBody string, config EmailConfig) error {
	m := mail.NewMessage()

	m.SetAddressHeader("From", config.SMTPUsername, config.SMTPSender)

	// 设置接收者
	m.SetHeader("To", receiver)
	// 设置邮件主题
	m.SetHeader("Subject", subject)
	// 设置邮件内容
	m.SetBody("text/html", htmlBody)

	d := mail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
