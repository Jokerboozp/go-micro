package main

import (
	"bytes"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

/**
该代码定义了一个 Mail 结构体，表示邮件配置信息。
结构体中包含了邮件服务器的相关字段。
该代码还定义了一个 Message 结构体，表示邮件消息的字段。
Mail 结构体包含了一个 SendSMTPMessage 方法，用于发送使用 SMTP 协议的邮件。
方法根据 Message 的字段构建邮件消息，并通过 SMTP 客户端发送邮件。
buildPlainTextMessage 方法根据模板构建纯文本格式的邮件消息。
buildHTMLMessage 方法根据模板构建 HTML 格式的邮件消息。
inlineCSS 方法将 CSS 样式嵌入 HTML 消息中。
getEncryption 方法根据加密方式的字符串值获取对应的枚举值。
通过这些方法，可以实现构建和发送邮件的功能。
*/

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	// 如果消息的发件人为空，则使用默认发件人地址
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	// 如果消息的发件人名称为空，则使用默认发件人名称
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// 创建数据映射，用于渲染邮件模板
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// 构建 HTML 格式的邮件消息
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	// 构建纯文本格式的邮件消息
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// 创建 SMTP 客户端
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// 创建邮件对象
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// 添加附件
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// 发送邮件
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

// 构建纯文本格式的邮件消息
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "/templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// 构建 HTML 格式的邮件消息
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "/templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	// 将 CSS 样式嵌入邮件消息中
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// 将 CSS 样式嵌入 HTML 中
func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

// 获取加密方式对应的枚举值
func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
