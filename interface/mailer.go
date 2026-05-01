package interfaces

import "github.com/michaelodeh/go-kit/pkg/dto"

type Mailer interface {
	SendMail(msg *dto.MailMessage) error
}

type MailTransport interface {
	Send(msg any) error
}
