package domain

import "errors"

var (
	ErrNotFound     = errors.New("water governance object not found")
	ErrConflict     = errors.New("water governance state conflict")
	ErrForbidden    = errors.New("water governance permission denied")
	ErrInvalidState = errors.New("invalid lifecycle transition")
	ErrExpired      = errors.New("session expired")
	ErrCancelled    = errors.New("operation cancelled")
	ErrDuplicate    = errors.New("duplicate idempotency key")
)
