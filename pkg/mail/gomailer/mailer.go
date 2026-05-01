package mailer

import (
	interfaces "github.com/alapa-ai/auth/internal/interface"
)

type goMailer struct {
	transport interfaces.MailTransport
}

func NewGoMailer(transport interfaces.MailTransport) interfaces.Mailer {
	return &goMailer{transport: transport}
}
