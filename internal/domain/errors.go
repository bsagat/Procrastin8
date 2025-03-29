package domain

import "errors"

var (
	ErrNotUniqueTask   = errors.New("task data must be unique")
	ErrActiveDateError = errors.New("active date field is invalid: time parsing error")
	ErrTaskNotFound    = errors.New("task is not found")
	ErrTaskChanged     = errors.New("task status is already changed")
	ErrInvalidStatus   = errors.New("status parameter is invalid,must be active or done")
)
