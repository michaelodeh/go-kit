package mailer

import interfaces "github.com/michaelodeh/go-kit/interface"

type goMailer struct {
	transport interfaces.MailTransport
}

func NewGoMailer(transport interfaces.MailTransport) interfaces.Mailer {
	return &goMailer{transport: transport}
}
