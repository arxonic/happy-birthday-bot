package auth

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/arxonic/gmh/internal/lib/logger/sl"
	"github.com/arxonic/gmh/internal/lib/token"
	"github.com/arxonic/gmh/internal/models"
)

var (
	ErrUserExists   = errors.New("user alredy exists")
	ErrUserNotFound = errors.New("user not found")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	userAuther   UserAuther
	emailSender  EmailSender
}

type EmailSender interface {
	SendEmail(to, subject, message string) error
}

type UserSaver interface {
	SaveAllUserInfo(models.User, models.UserMessenger, models.Organization) (int64, error)
}

type UserProvider interface {
	IsActivated(messengerType string, messengerID, chatID int64) (bool, error)
	// UserByEmail(email string) (models.User, error)
	// UserMessenger(messenger string, messengerID int64) (models.UserMessenger, error)
}

type UserAuther interface {
	UpdateUserActivationStatus(messengerType string, messengerID, chatID int64, token string) error
}

// New returns a new instance of the Auth service.
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, userAuther UserAuther, noty EmailSender) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		userAuther:   userAuther,
		emailSender:  noty,
	}
}

func (a *Auth) RegisterNewUser(user models.User, userMessenger models.UserMessenger, userOrganization models.Organization) (int64, error) {
	const fn = "auth.Register"

	log := a.log.With(slog.String("fn", fn))

	// Generate auth token
	token, err := token.NewToken()
	if err != nil {
		log.Error("failed to generate auth token", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	userMessenger.Token = token

	// Send Email
	authLink := fmt.Sprintf(
		"http://localhost:2001/v1/auth?token=%s&mtype=%s&mid=%d&chatid=%d&redirect=%s",
		token,
		userMessenger.MessengerType,
		userMessenger.MessengerID,
		userMessenger.ChatID,
		"https://t.me/GPMHappyBBot",
	)

	err = a.emailSender.SendEmail(user.Email, "Регистрация", authLink)
	if err != nil {
		log.Error("failed to send email", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	// Save User
	uID, err := a.userSaver.SaveAllUserInfo(user, userMessenger, userOrganization)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return uID, nil
}

// IsActivated return Activation Account Status if UserMessenger exists
func (a *Auth) IsActivated(messengerType string, messengerID, chatID int64) (bool, error) {
	const fn = "auth.IsActivated"

	log := a.log.With(slog.String("fn", fn))

	isActivated, err := a.userProvider.IsActivated(messengerType, messengerID, chatID)
	if err != nil {
		log.Debug("rows not found", sl.Err(err))
		return false, err
	}

	return isActivated, nil
}

// AccountActivation change Users Account Activation status if he followed the auth link
func (a *Auth) AccountActivation(messengerType string, messengerID, chatID int64, token string) error {
	const fn = "auth.AccountActivation"

	log := a.log.With(slog.String("fn", fn))

	if err := a.userAuther.UpdateUserActivationStatus(messengerType, messengerID, chatID, token); err != nil {
		log.Error("failed to save user activation status", sl.Err(err))
		return err
	}

	return nil
}
