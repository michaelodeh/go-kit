package utils

import "fmt"

func SendMail(email string, subject string, body string) error {
	fmt.Println("Sending mail to", email)
	fmt.Println("Subject:", subject)
	fmt.Println("Body:", body)
	return nil
}
