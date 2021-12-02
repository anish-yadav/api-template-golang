package mail

import (
	"fmt"
	"net/smtp"
)

var auth smtp.Auth

func Init() {
	// Authentication.
	auth = smtp.PlainAuth("", "", "password", "smtp.gmail.com")
}

func SendEmail(to []string) {
	message := []byte("This is a test email message.")



	// Sending email.
	err := smtp.SendMail("smtp.gmail.com"+":"+"567", auth, "me@mail.com", to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
}