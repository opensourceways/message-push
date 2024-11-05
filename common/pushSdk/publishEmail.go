package pushSdk

import (
	"bytes"
	"crypto/tls"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
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

func mdToHtml(body string) (string, error) {
	// 创建一个缓冲区以写入转换后的 HTML
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote),
		goldmark.WithRendererOptions(html.WithUnsafe()))
	err := md.Convert([]byte(body), &buf)
	if err != nil {
		return "", err
	}

	// 返回转换后的 HTML 字符串
	return buf.String(), nil
}

func sendSSLEmail(receiver, subject, body string, config EmailConfig) error {
	m := mail.NewMessage()

	m.SetAddressHeader("From", config.SMTPUsername, config.SMTPSender)

	// 设置接收者
	m.SetHeader("To", receiver)
	// 设置邮件主题
	m.SetHeader("Subject", subject)
	// 设置邮件内容
	htmlBody, err := mdToHtml(body)
	if err != nil {
		return err
	}
	logrus.Infof("the data is %v\n, the html is %v", body, htmlBody)
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
