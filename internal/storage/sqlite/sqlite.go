package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arxonic/gmh/internal/models"
	repo "github.com/arxonic/gmh/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	return &Storage{db: db}, nil
}

// User return User model by UserID
func (s *Storage) User(uID int64) (models.User, error) {
	const fn = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, first_name, last_name, patronymic, birth_date, email FROM users WHERE id = ?")
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(uID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Patronymic, &user.BirthDate, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repo.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s:%w", fn, err)
	}

	return user, nil
}

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

func (s *Storage) UserByEmail(email string) (models.User, error) {
	const fn = "storage.sqlite.UserByEmail"

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(email).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repo.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s:%w", fn, err)
	}

	return user, nil
}

// SaveAllUserInfo save User into storage and return his UserID
func (s *Storage) SaveAllUserInfo(
	user models.User,
	userMessenger models.UserMessenger,
	userOrganization models.Organization,
) (int64, error) {
	const fn = "storage.sqlite.SaveAllUserInfo"

	uID, err := s.SaveUser(user)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	var orgID int64

	// TODO FIX THIS (IF OR EXISTS, SEARCH ORG BY FULL NAME)
	// Save organization if not exist
	orgID, err = s.SaveOrganization(userOrganization)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	_, err = s.SaveUserOrganization(uID, orgID)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	userMessenger.UserID = uID
	_, err = s.SaveUserMessenger(userMessenger)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return uID, nil
}

// SaveUser save User into storage and return his UserID
func (s *Storage) SaveUser(user models.User) (int64, error) {
	const fn = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (first_name, last_name, patronymic, birth_date, email) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.FirstName, user.LastName, user.Patronymic, user.BirthDate, user.Email)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s:%w", fn, repo.ErrUserExists)
		}

		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return id, nil
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

// Organization return Organization model by OrganizationID
func (s *Storage) Organization(orgID int64) (models.Organization, error) {
	const fn = "storage.sqlite.Organization"

	stmt, err := s.db.Prepare("SELECT id, name, city, office, department FROM organizations WHERE id = ?")
	if err != nil {
		return models.Organization{}, err
	}
	defer stmt.Close()

	var org models.Organization
	err = stmt.QueryRow(orgID).Scan(&org.ID, &org.Name, &org.City, &org.Office, &org.Department)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Organization{}, repo.ErrOrganizationNotFound
		}
		return models.Organization{}, fmt.Errorf("%s:%w", fn, err)
	}

	return org, nil
}

// SaveOrganization save Organization into storage and return his OrganizationID
func (s *Storage) SaveOrganization(org models.Organization) (int64, error) {
	const fn = "storage.sqlite.SaveOrganization"

	// Check id Exist
	stmt, err := s.db.Prepare("SELECT id FROM organizations WHERE name = ? AND city = ? AND office = ? AND department = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(org.Name, org.City, org.Office, org.Department).Scan(&id)
	if !errors.Is(err, sql.ErrNoRows) {
		if !errors.Is(err, sql.ErrNoRows) {
			return id, nil
		}
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	// Save
	stmt, err = s.db.Prepare("INSERT INTO organizations (name, city, office, department) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(org.Name, org.City, org.Office, org.Department)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s:%w", fn, repo.ErrOrganizationExists)
		}

		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return id, nil
}

// SaveUserOrganization save assignments by UserID and OrganizationID and return new row ID
func (s *Storage) SaveUserOrganization(uID, orgID int64) (int64, error) {
	const fn = "storage.sqlite.SaveUserOrganization"

	stmt, err := s.db.Prepare("INSERT INTO user_organizations (user_id, organization_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(uID, orgID)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return id, nil
}

// FindOrgByFields return strings DISTINCT rows by fields. For example args: Gazprom, Moscow returns all DISTINCT offices of this org.
// !!! IF YOU PASSED ALL FIELDS IN THE FUNC - RETURNS []ID `organization` table
func (s *Storage) FindOrgByFields(fields ...string) ([]string, error) {
	const fn = "storage.sqlite.FindOrgByFields"

	q := ""
	switch len(fields) {
	case 0:
		q = "SELECT DISTINCT name FROM organizations"
	case 1:
		q += "SELECT DISTINCT city FROM organizations WHERE name = ?"
	case 2:
		q += "SELECT DISTINCT office FROM organizations WHERE name = ? AND city = ?"
	case 3:
		q += "SELECT DISTINCT department FROM organizations WHERE name = ? AND city = ? AND office = ?"
	case 4:
		q += "SELECT DISTINCT id FROM organizations WHERE name = ? AND city = ? AND office = ? AND department = ?"
	}

	stmt, err := s.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	switch len(fields) {
	case 0:
		rows, err = stmt.Query()
	case 1:
		rows, err = stmt.Query(fields[0])
	case 2:
		rows, err = stmt.Query(fields[0], fields[1])
	case 3:
		rows, err = stmt.Query(fields[0], fields[1], fields[2])
	case 4:
		rows, err = stmt.Query(fields[0], fields[1], fields[2], fields[3])
	}
	if err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	orgs := make([]string, 0)
	for rows.Next() {
		var org string
		if err := rows.Scan(&org); err != nil {
			return nil, fmt.Errorf("%s:%w", fn, err)
		}

		orgs = append(orgs, org)
	}

	return orgs, nil
}

func (s *Storage) UserIDsByOrgID(id int64) ([]int64, error) {
	const fn = "storage.sqlite.UserIDsByOrgID"

	stmt, err := s.db.Prepare("SELECT user_id FROM user_organizations WHERE organization_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("%s:%w", fn, err)
		}

		ids = append(ids, id)
	}

	return ids, nil
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

func (s *Storage) Subscribe(subID, uID int64) (int64, error) {
	const fn = "storage.sqlite.Subscribe"

	stmt, err := s.db.Prepare("INSERT INTO subscribes (user_id, sub_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(uID, subID)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	return id, nil
}
