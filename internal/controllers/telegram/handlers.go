package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arxonic/gmh/internal/controllers/telegram/states"
	"github.com/arxonic/gmh/internal/lib/email"
	"github.com/arxonic/gmh/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		b.SendMessage(m, "Давайте поищем Ваших коллег, на чьи дни рожения Вы хотите подписаться")

		orgs, err := uf.FindUser()
		if err != nil {
			b.SendMessage(m, "Такого варианта нет")
			return states.StateMenu, nil
		}

		tgbotapi.NewReplyKeyboard()
		rows := make([]tgbotapi.KeyboardButton, 0)
		for _, org := range orgs {
			r := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(org))
			rows = append(rows, r...)
		}

		markup := tgbotapi.NewReplyKeyboard(rows)

		b.SendKeyboardMessage(m, "Выберите оргагизацию из списка (введите её название):", markup)

		return states.StateFind, nil

	case "2":
		b.SendMessage(m, "Раздел в разработке...")
		return states.StateMenu, nil
	default:
		b.handleUnknownCommand(m)
		return states.StateMenu, nil
	}
}

type UserFinder interface {
	User(id int64) (models.User, error)
	UsersByOrgID(id int64) ([]models.User, error)
	FindUser(...string) ([]string, error)
	Subscribe(messengerID, uID int64) error
}

func (b *Bot) FinderHandler(m *tgbotapi.Message, uf UserFinder, state *states.UserState) (int, error) {
	response := ""

	orgs := make([]string, 0)

	var err error

	if state.Finder.Organization == "" {
		orgName := m.Text

		orgs, err = uf.FindUser(orgName)
		if err != nil {
			b.SendMessage(m, "Такого варианта нет")
			return states.StateMenu, nil
		}

		state.Finder.Organization = orgName

		response = "Выберите город этой организации из списка:"

	} else if state.Finder.City == "" {
		orgCity := m.Text

		orgs, err = uf.FindUser(state.Finder.Organization, orgCity)
		if err != nil {
			b.SendMessage(m, "Такого варианта нет")
			return states.StateMenu, nil
		}

		state.Finder.City = orgCity

		response = "Выберите офис этой организации из списка:"

	} else if state.Finder.Office == "" {
		orgOffice := m.Text

		orgs, err = uf.FindUser(state.Finder.Organization, state.Finder.City, orgOffice)
		if err != nil {
			b.SendMessage(m, "Такого варианта нет")
			return states.StateMenu, nil
		}

		state.Finder.Office = orgOffice

		response = "Выберите отдел этой организации из списка:"

	} else if state.Finder.Department == "" {
		orgDepart := m.Text

		ids, err := uf.FindUser(state.Finder.Organization, state.Finder.City, state.Finder.Office, orgDepart)
		if err != nil {
			b.SendMessage(m, "Такого варианта нет")
			return states.StateMenu, nil
		}

		id, err := strconv.Atoi(ids[0])
		if err != nil {
			b.SendMessage(m, "Ошибка сервера, повторите попытку позже")
			return states.StateMenu, nil
		}

		users, err := uf.UsersByOrgID(int64(id))
		if err != nil {
			b.SendMessage(m, "Ошибка сервера, повторите попытку позже")
			return states.StateMenu, nil
		}

		state.Finder.Department = orgDepart

		for _, u := range users {
			r := fmt.Sprintf("%d %s %s %s", u.ID, u.LastName, u.FirstName, u.LastName)
			orgs = append(orgs, r)
		}

		response = "Выберите человека из списка или введите его id:"
	} else if state.Finder.Office == "" {
		orgOffice := m.Text

		orgs, err = uf.FindUser(state.Finder.Organization, state.Finder.City, orgOffice)
		if err != nil {
			b.SendMessage(m, "Такого варианта нет")
			return states.StateMenu, nil
		}

		state.Finder.Office = orgOffice

		response = "Выберите отдел этой организации из списка:"

	} else if state.Finder.UserID == 0 {
		usr := strings.Split(m.Text, " ")

		usrID, err := strconv.Atoi(usr[0])
		if err != nil {
			b.SendMessage(m, "Ошибка ввода")
			return states.StateMenu, nil
		}

		err = uf.Subscribe(m.From.ID, int64(usrID))
		if err != nil {
			b.SendMessage(m, "Ошибка сервера, повторите попытку позже")
			return states.StateMenu, nil
		}

		// null state
		state.Finder = states.FindState{}

		b.SendMessage(m, "Вы успешно подписались на польщователя! Я отправлю вам ссылку на чат за неделю до дня рождения вашего коллеги <3")

		return states.StateMenu, nil
	}

	// Draw buttons
	markup := createReplyKeyboardMarkup(orgs)

	b.SendKeyboardMessage(m, response, markup)

	return state.State, nil
}

func (b *Bot) handleUnknownCommand(m *tgbotapi.Message) error {
	text := "Эта команда мне незнакома 0_o"
	return b.SendMessage(m, text)
}

// createReplyKeyboardMarkup create buttons from []string
func createReplyKeyboardMarkup(btns []string) tgbotapi.ReplyKeyboardMarkup {
	rows := make([]tgbotapi.KeyboardButton, 0)
	for _, b := range btns {
		r := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(b))
		rows = append(rows, r...)
	}

	return tgbotapi.NewReplyKeyboard(rows)
}
