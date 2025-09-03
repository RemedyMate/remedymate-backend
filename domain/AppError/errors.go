package AppError

import "errors"

// sentinel errors for the application
var (
	ErrTopicNotFound        = errors.New("topic not found")
	ErrLanguageNotAvailable = errors.New("language not available")
	ErrUnsupportedLanguage  = errors.New("unsupported language")
	ErrNoTopicMapped        = errors.New("no topic could be mapped from the provided symptoms")

	// user errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserStatusNotFound = errors.New("user status not found")
	ErrEmailAlreadyExist  = errors.New("email already exist")
	ErrIncorrectPassword  = errors.New("incorrect password")
	ErrInactiveAccount    = errors.New("account is inactive")
	ErrInvalidToken       = errors.New("invalid token")

	// server errors
	ErrInternalServer = errors.New("internal server error")
)
