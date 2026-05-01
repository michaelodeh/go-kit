package mailer

import (
	interfaces "github.com/alapa-ai/auth/internal/interface"
	"github.com/alapa-ai/auth/pkg/dto"
	"gopkg.in/gomail.v2"
)

type GoMailerDialer struct {
	dialer *gomail.Dialer
}

func NewGoMailerDialer(config *dto.MailConfig) interfaces.MailTransport {
	return &GoMailerDialer{
		dialer: gomail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
		),
	}

}

func (t *GoMailerDialer) Send(msg any) error {
	return t.dialer.DialAndSend(msg.(*gomail.Message))
}
