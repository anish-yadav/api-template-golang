package util

import "net/mail"

func ValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
