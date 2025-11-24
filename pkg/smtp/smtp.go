package smtp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

// TencentSendToMailFunc 腾讯云 特殊处理
func TencentSendToMailFunc(identity, user, password string, host string, port int) func(opts ...SendOption) error {
	auth := smtp.PlainAuth(identity, user, password, host)
	return func(opts ...SendOption) error {
		option := NewSendOptions(opts...)
		if option.From == "" {
			option.From = user
		}
		if len(option.To) > 1 {
			return errors.New("腾讯云邮件推送暂不支持多目标")
		}
		return SendMailWithTLS(fmt.Sprintf("%s:%d", host, port), auth, user, option.Targets(), []byte(option.Encode()))

	}
}

func SendToMailFunc(identity, user, password string, host string, port int) func(opts ...SendOption) error {
	auth := smtp.PlainAuth(identity, user, password, host)
	return func(opts ...SendOption) error {
		option := NewSendOptions(opts...)
		if option.From == "" {
			option.From = user
		}
		return smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, user, option.Targets(), []byte(option.Encode()))
	}
}

func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("tls.Dial Error:", err)
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// SendMailWithTLS send email with tls
func SendMailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smtp client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
