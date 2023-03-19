package utility

import "net/mail"

func CheckIsValidMail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
