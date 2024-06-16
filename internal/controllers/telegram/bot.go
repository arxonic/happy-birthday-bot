package telegram

import (
	"log/slog"

	"github.com/arxonic/gmh/internal/controllers/telegram/states"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MessengerType = "telegram"

type Bot struct {
	bot *tgbotapi.BotAPI
	log *slog.Logger
}

func NewBot(tgBotKey string, log *slog.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(tgBotKey)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot: bot,
		log: log,
	}, nil
}

func (b *Bot) Run(states *states.States, uf UserFinder, ua UserAuther, emp Employer) {
	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60 // TODO save to config

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		m := update.Message
		if m == nil {
			continue
		}

		b.handleState(m, states, uf, ua, emp)
	}
}

func (b *Bot) handleState(m *tgbotapi.Message, s *states.States, uf UserFinder, ua UserAuther, emp Employer) {
	userID := m.From.ID

	// Инициализация состояния пользователя, если его еще нет
	if _, ok := s.UserStates[userID]; !ok {
		s.UserStates[userID] = &states.UserState{State: states.StateAuthMiddleware}
	}

	state := s.UserStates[userID]

	switch state.State {
	case states.StateAuthMiddleware:
		newState, err := b.Auth(m, ua)
		if err != nil {
			return
		}
		s.UserStates[userID].State = newState

	case states.StateEmailWait:
		newState, err := b.EmailWait(m, emp, ua)
		if err != nil {
			return
		}
		s.UserStates[userID].State = newState

	case states.StateMenu:
		newState, err := b.MenuHandler(m, uf)
		if err != nil {
			return
		}
		s.UserStates[userID].State = newState

	case states.StateOrgName:
		newState, err := b.FinderHandler(m, uf, state)
		if err != nil {
			return
		}
		s.UserStates[userID].State = newState
	}

}

func (b *Bot) SendMessage(m *tgbotapi.Message, text string) error {
	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ReplyToMessageID = m.MessageID

	_, err := b.bot.Send(msg)

	return err
}
