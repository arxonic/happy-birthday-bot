package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arxonic/gmh/internal/models"
	repo "github.com/arxonic/gmh/internal/storage"
)

// UpdateUserActivationStatus is activating user account
func (s *Storage) UpdateUserActivationStatus(messengerType string, messengerID, chatID int64, token string) error {
	const fn = "storage.sqlite.SaveUserActivationStatus"

	q := `UPDATE user_messengers SET is_activated = ?, token = ? 
	WHERE messenger_type = ? AND messenger_id = ? AND chat_id = ? AND token = ?`
	_, err := s.db.Exec(
		q,
		1,
		"",
		messengerType,
		messengerID,
		chatID,
		token,
	)
	if err != nil {
		return fmt.Errorf("%s:%w", fn, err)
	}

	return nil
}

// IsActivated return Activation Account status by Messenger Info
func (s *Storage) IsActivated(messengerType string, messengerID, chatID int64) (bool, error) {
	const fn = "storage.sqlite.IsActivated"

	stmt, err := s.db.Prepare("SELECT is_activated FROM user_messengers WHERE messenger_type = ? AND messenger_id = ? AND chat_id = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var isActivated bool
	err = stmt.QueryRow(messengerType, messengerID, chatID).Scan(&isActivated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repo.ErrUserNotFound
		}
		return false, fmt.Errorf("%s:%w", fn, err)
	}

	return isActivated, nil
}

// UserMessenger return UserMessenger model by Messenger Type and users MessengerID into this messenger
func (s *Storage) UserMessenger(messenger string, messengerID int64) (models.UserMessenger, error) {
	const fn = "storage.sqlite.UserMessenger"

	stmt, err := s.db.Prepare("SELECT * FROM user_messengers WHERE messenger_type = ? AND messenger_id = ?")
	if err != nil {
		return models.UserMessenger{}, err
	}
	defer stmt.Close()

	var userMessenger models.UserMessenger
	err = stmt.QueryRow(messenger, messengerID).Scan(&userMessenger)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserMessenger{}, repo.ErrUserNotFound
		}
		return models.UserMessenger{}, fmt.Errorf("%s:%w", fn, err)
	}

	return userMessenger, nil
}

// SaveUserMessenger save user messenger info and return new row ID
func (s *Storage) SaveUserMessenger(data models.UserMessenger) (int64, error) {
	const fn = "storage.sqlite.SaveUserMessenger"

	stmt, err := s.db.Prepare("INSERT INTO user_messengers (user_id, messenger_type, messenger_id, chat_id, is_activated, token) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		data.UserID,
		data.MessengerType,
		data.MessengerID,
		data.ChatID,
		data.IsActivated,
		data.Token,
	)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return id, nil
}

func (s *Storage) UserIDByMessengerID(id int64) (int64, error) {
	const fn = "storage.sqlite.UserIDByMessengerID"

	stmt, err := s.db.Prepare("SELECT user_id FROM user_messengers WHERE messenger_id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var uID int64
	err = stmt.QueryRow(id).Scan(&uID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repo.ErrOrganizationNotFound
		}
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return uID, nil
}
