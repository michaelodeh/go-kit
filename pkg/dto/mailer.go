package dto

type MailMessage struct {
	From        string
	To          string
	Cc          string
	Subject     string
	TextBody    string
	HTMLBody    string
	Attachments []string
}

type MailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	UseTLSOrSSL bool
}
