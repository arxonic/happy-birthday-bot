package telegram

import (
	"strconv"

	"github.com/arxonic/gmh/internal/controllers/telegram/states"
	"github.com/arxonic/gmh/internal/lib/email"
	"github.com/arxonic/gmh/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserProvider interface {
}

type UserAuther interface {
	RegisterNewUser(models.User, models.UserMessenger, models.Organization) (int64, error)
	IsActivated(messengerType string, messengerID, chatID int64) (bool, error)
}

func (b *Bot) Auth(m *tgbotapi.Message, ua UserAuther) (int, error) {
	isActivated, err := ua.IsActivated(MessengerType, m.From.ID, m.Chat.ID)
	if err != nil {
		// if New user
		text := "Добри день! Вас приветствует тг бот для поздравления сотрудников Газпром-медиа с Днем Рождения! Введите свою корпоративную почту:"
		b.SendMessage(m, text)
		return states.StateEmailWait, nil
	}
	if isActivated {
		b.SendMessage(m, "Здравствуйте!")
		b.SendMessage(m, "Введите: \n(1) - Найти коллегу, чтобы подписаться на его ДР\n(2) - Список ваших подписок")
		return states.StateMenu, nil
	} else {
		// if user not follow auth link
		b.SendMessage(m, "На указанную почту пришло письмо со ссылкой, перейдите по ней :з")
		return states.StateAuthMiddleware, nil
	}
}

type Employer interface {
	Employee(email string) (models.Emp, error)
}

func (b *Bot) EmailWait(m *tgbotapi.Message, emp Employer, ua UserAuther) (int, error) {
	e := m.Text
	if !email.Valid(e) {
		b.SendMessage(m, "Неверный формат почты. Давайте попробуем еще раз. Введите свою корпоративную почту:")
		return states.StateEmailWait, nil
	}

	// Employer API request
	employee, err := emp.Employee(e)
	if err != nil {
		b.SendMessage(m, "Мы не нашли Вашу почту в организации. Давайте попробуем еще раз. Введите свою корпоративную почту:")
		return states.StateEmailWait, nil
	}

	// Convert Employer response to User model
	user, org := models.EmpToUser(employee)
	mess := models.UserMessenger{
		MessengerType: MessengerType,
		MessengerID:   m.From.ID,
		ChatID:        m.Chat.ID,
	}

	// Save user
	_, err = ua.RegisterNewUser(user, mess, org)
	if err != nil {
		b.SendMessage(m, "Возникла ошибка при обработке данных. Повторите попытку позже")
		return states.StateEmailWait, nil
	}

	b.SendMessage(m, "На указанную почту пришло письмо со ссылкой, перейдите по ней :з")

	return states.StateAuthMiddleware, nil
}

func (b *Bot) MenuHandler(m *tgbotapi.Message, uf UserFinder) (int, error) {
	switch m.Text {
	case "1":
		b.SendMessage(m, "Давайте поищем Ваших коллег, на чьи дни рожения Вы хотите подписаться. Выберите оргагизацию из списка (введите её название):")
		orgs, err := uf.FindUser()
		if err != nil {
			b.SendMessage(m, "Такого варианта нет")
			return states.StateMenu, nil
		}

		// TODO change string concatinate
		resp := ""
		for i, org := range orgs {
			resp += strconv.Itoa(i) + " - " + org.Name + "\n"
		}

		b.SendMessage(m, resp)
		return states.StateOrgName, nil
	case "2":
		b.SendMessage(m, "Раздел в разработке...")
		return states.StateMenu, nil
	default:
		b.handleUnknownCommand(m)
		return states.StateMenu, nil
	}
}

type UserFinder interface {
	FindUser(...string) ([]models.Organization, error)
}

func (b *Bot) FinderHandler(m *tgbotapi.Message, uf UserFinder, state *states.UserState) (int, error) {
	switch state.State {
	case states.StateOrgName:
		orgName := m.Text

		orgs, err := uf.FindUser(orgName)
		if err != nil {
			return state.State, nil
		}

		state.Finder.Organization = orgName

		// TODO change string concatinate
		resp := ""
		for i, org := range orgs {
			resp += strconv.Itoa(i) + org.City + "\n"
		}

		b.SendMessage(m, resp)

		return states.StateOrgCity, nil

	case states.StateOrgCity:
		// orgCity := m.Text

		// orgs, err := uf.FindOrgByFields(orgName, orgCity)
		// if err != nil {
		// 	return state.State, nil
		// }

		return states.StateOrgCity, nil

	case states.StateOrgOffice:
		// orgName := m.Text

		// orgs, err := uf.FindOrgByFields(orgName)
		// if err != nil {
		// 	return state.State, nil
		// }

		return states.StateOrgCity, nil

	case states.StateOrgDepart:
		// orgName := m.Text

		// orgs, err := uf.FindOrgByFields(orgName)
		// if err != nil {
		// 	return state.State, nil
		// }

		return states.StateSubscribe, nil

	default:
		return state.State, nil

	}
}

func (b *Bot) handleUnknownCommand(m *tgbotapi.Message) error {
	text := "Эта команда мне незнакома 0_o"
	return b.SendMessage(m, text)
}
