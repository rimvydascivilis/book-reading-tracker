package domain

import (
	"errors"
)

var (
	ErrValidation     = errors.New("validation error")
	ErrAuthentication = errors.New("authentication error")
	ErrRecordNotFound = errors.New("record not found")
)
