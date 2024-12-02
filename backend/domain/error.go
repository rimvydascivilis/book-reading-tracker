package domain

import (
	"errors"
)

var (
	ErrValidation     = errors.New("validation error")
	ErrAuthentication = errors.New("authentication error")
	ErrRecordNotFound = errors.New("record not found")
	ErrForbidden      = errors.New("forbidden")
	ErrAlreadyExists  = errors.New("record already exists")
)
