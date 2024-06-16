package subscribe

import (
	"log/slog"

	"github.com/arxonic/gmh/internal/models"
)

type Sub struct {
	log *slog.Logger
	Subscriber
}

type Subscriber interface {
	FindOrgByFields(...string) ([]models.Organization, error)
	UsersByOrgID(id int64) ([]models.User, error)
	UserIDByMessengerID(id int64) (int64, error)
	Subscribe(subID, uID int64) (int64, error)
}

func New(log *slog.Logger, s Subscriber) *Sub {
	return &Sub{
		log:        log,
		Subscriber: s,
	}
}

func (s *Sub) FindUser(fields ...string) ([]models.Organization, error) {
	orgs, err := s.FindOrgByFields(fields...)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
