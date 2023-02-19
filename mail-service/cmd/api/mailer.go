package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

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
	From       string
	FromName   string
	To         string
	Subject    string
	Attachment []string
	Data       any
	DataMap    map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	/*
	 * email could be send in various format, e.g. plain text, html
	 */

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return nil
	}

	plainTextMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return nil
	}

	mailServer := mail.NewSMTPClient()
	mailServer.Host = m.Host
	mailServer.Port = m.Port
	mailServer.Username = m.Username
	mailServer.Password = m.Password
	mailServer.Encryption = m.getEncryption(m.Encryption)
	mailServer.KeepAlive = false
	mailServer.ConnectTimeout = 10 * time.Second
	mailServer.SendTimeout = 10 * time.Second

	smtpClient, err := mailServer.Connect()
	if err != nil {
		return nil
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainTextMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachment) > 0 {
		for _, x := range msg.Attachment {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil

}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", nil
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", nil
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", nil
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", nil
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", nil
	}

	plainTextMessage := tpl.String()

	return plainTextMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", nil
	}

	html, err := prem.Transform()
	if err != nil {
		return "", nil
	}

	return html, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
