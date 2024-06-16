package storage

import (
	"errors"
)

var (
	ErrUserExists           = errors.New("user alredy exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrOrganizationExists   = errors.New("organization alredy exists")
	ErrOrganizationNotFound = errors.New("organization not found")
)
