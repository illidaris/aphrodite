package smtp

import (
	"net/smtp"
)

func SendToMailFunc(identity, user, password string, host string) func(opts ...SendOption) error {
	auth := smtp.PlainAuth(identity, user, password, host)
	return func(opts ...SendOption) error {
		option := NewSendOptions(opts...)
		if option.From == "" {
			option.From = user
		}
		return smtp.SendMail(host, auth, user, option.Targets(), []byte(option.Encode()))
	}
}
