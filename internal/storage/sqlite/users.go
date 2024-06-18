package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arxonic/gmh/internal/models"
	repo "github.com/arxonic/gmh/internal/storage"
	"github.com/mattn/go-sqlite3"
)

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

func (s *Storage) UserByEmail(email string) (models.User, error) {
	const fn = "storage.sqlite.UserByEmail"

	stmt, err := s.db.Prepare("SELECT id, first_name, last_name, patronymic, birth_date, email FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Patronymic, &user.BirthDate, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repo.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s:%w", fn, err)
	}

	return user, nil
}

// UsersWhoseBirthdayIsInXDays return []User whose birthday is in X days
func (s *Storage) UsersWhoseBirthdayIsInXDays(x int) ([]models.User, error) {
	const fn = "storage.sqlite.UsersWhoseBirthdayIsInXDays"

	stmt, err := s.db.Prepare("SELECT id, first_name, last_name, patronymic, birth_date, email " +
		"FROM users WHERE DATE_FORMAT(birthdate, '%m-%d') " +
		"BETWEEN DATE_FORMAT(DATE_ADD(CURDATE(), INTERVAL 1 DAY), '%m-%d') " +
		"AND DATE_FORMAT(DATE_ADD(CURDATE(), INTERVAL ? DAY), '%m-%d')")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(x)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Patronymic, &user.BirthDate, &user.Email); err != nil {
			return nil, fmt.Errorf("%s:%w", fn, err)
		}

		users = append(users, user)
	}

	return users, nil
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
