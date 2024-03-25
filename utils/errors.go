package utils

import (
  "errors"
)

var (
	ErrInvalidEmail       = errors.New("Invalid email")
	ErrEmailAlreadyExists = errors.New("Email already exists")
	ErrEmptyPassword      = errors.New("Password can't be empty")
	ErrInvalidAuthToken   = errors.New("Invalid auth-token")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrUnauthorized       = errors.New("Unauthorized")
)
