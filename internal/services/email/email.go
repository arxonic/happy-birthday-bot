package email

import (
	"log/slog"
)

type Sender struct {
	log *slog.Logger
	EmailSender
}

type EmailSender interface {
	SendEmail(to, subject, message string) error
}

func New(log *slog.Logger, sender EmailSender) *Sender {
	return &Sender{
		log:         log,
		EmailSender: sender,
	}
}
