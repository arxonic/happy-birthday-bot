package subscribe

import (
	"log/slog"

	"github.com/arxonic/gmh/internal/models"
)

type Sub struct {
	log *slog.Logger
	Subscriber
	UserProvider
}

type UserProvider interface {
	User(id int64) (models.User, error)
	UserIDsByOrgID(id int64) ([]int64, error)
	UserIDByMessengerID(id int64) (int64, error)
}

type Subscriber interface {
	FindOrgByFields(...string) ([]string, error)
	Subscribe(subID, uID int64) (int64, error)
}

func New(log *slog.Logger, s Subscriber, up UserProvider) *Sub {
	return &Sub{
		log:          log,
		Subscriber:   s,
		UserProvider: up,
	}
}

func (s *Sub) Subscribe(messengerSubID, uID int64) error {
	subID, err := s.UserProvider.UserIDByMessengerID(messengerSubID)
	if err != nil {
		return err
	}

	_, err = s.Subscriber.Subscribe(subID, uID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sub) FindUser(fields ...string) ([]string, error) {
	orgs, err := s.FindOrgByFields(fields...)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (s *Sub) User(id int64) (models.User, error) {
	user, err := s.UserProvider.User(id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Sub) UsersByOrgID(id int64) ([]models.User, error) {
	usersIDs, err := s.UserProvider.UserIDsByOrgID(id)
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)

	for _, id := range usersIDs {
		user, err := s.UserProvider.User(id)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
