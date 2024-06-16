package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Sender struct {
	d *gomail.Dialer
}

func New(host string, port int, sender string, password string) *Sender {
	d := gomail.NewDialer(host, port, sender, password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Sender{
		d: d,
	}
}

func (s *Sender) SendEmail(to, subject, message string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.d.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", message)
	if err := s.d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
