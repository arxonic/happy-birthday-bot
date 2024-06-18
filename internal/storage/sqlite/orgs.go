package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arxonic/gmh/internal/models"
	repo "github.com/arxonic/gmh/internal/storage"
	"github.com/mattn/go-sqlite3"
)

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

	// Check if Exist
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
