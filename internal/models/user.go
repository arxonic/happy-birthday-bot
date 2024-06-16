package models

import "time"

// User представляет информацию о пользователе
type User struct {
	ID         int64     `db:"id" json:"uid"`
	FirstName  string    `db:"first_name" json:"first_name"`
	LastName   string    `db:"last_name" json:"last_name"`
	Patronymic string    `db:"patronymic" json:"patronymic"`
	BirthDate  time.Time `db:"birth_date" json:"birth_date"`
	Email      string    `db:"email" json:"email"`
}

// Organization представляет информацию об организации
type Organization struct {
	ID         int64  `db:"id" json:"org_id"`
	Name       string `db:"name" json:"name"`
	City       string `db:"city" json:"city"`
	Office     string `db:"office" json:"office"`
	Department string `db:"department" json:"department"`
}

// UserOrganization связывает пользователя с организацией
type UserOrganization struct {
	UserID         int64 `db:"user_id" json:"user_id"`
	OrganizationID int64 `db:"organization_id" json:"organization_id"`
}

// UserMessenger представляет информацию о мессенджере пользователя
type UserMessenger struct {
	UserID        int64  `db:"user_id" json:"user_id"`
	MessengerType string `db:"messenger_type" json:"messenger_type"`
	MessengerID   int64  `db:"messenger_id" json:"messenger_id"`
	ChatID        int64  `db:"chat_id" json:"chat_id"`
	IsActivated   bool   `db:"is_activated" json:"is_activated"`
	Token         string `db:"token" json:"token"`
}
