package mailer

import (
	"github.com/michaelodeh/go-kit/dto"
	"gopkg.in/gomail.v2"
)

func (m *goMailer) SendMail(message *dto.MailMessage) error {
	msg := buildMessage(message)
	return m.transport.Send(msg)
}

func buildMessage(m *dto.MailMessage) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.To)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.HTMLBody)
	return msg
}
