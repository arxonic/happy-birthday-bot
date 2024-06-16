package v1

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type UserAuther interface {
	AccountActivation(messengerType string, messengerID, chatID int64, token string) error
}

func NewRouts(handler *chi.Mux, log *slog.Logger, auther UserAuther) {
	handler.Get("/v1/auth", Auth(log, auther))
}
