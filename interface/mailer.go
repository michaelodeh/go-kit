package interfaces

import "github.com/michaelodeh/go-kit/dto"

type Mailer interface {
	SendMail(msg *dto.MailMessage) error
}

type MailTransport interface {
	Send(msg any) error
}
