package usecase

import "errors"

var (
	InvalidTokenErr = errors.New("Token invalid")
	ExpiredTokenErr = errors.New("Token expired")
	SessionNotFound = errors.New("Session not found")
)
