package AppError

import "errors"

// sentinel errors for the application
var (
	ErrTopicNotFound        = errors.New("topic not found")
	ErrLanguageNotAvailable = errors.New("language not available")
	ErrUnsupportedLanguage  = errors.New("unsupported language")
	ErrNoTopicMapped        = errors.New("no topic could be mapped from the provided symptoms")
	ErrInvalidInput         = errors.New("invalid input")
	ErrTopicAlreadyExists   = errors.New("topic already exists")

	// user errors
	ErrUserNotFound         = errors.New("user not found")
	ErrUserNotAuthenticated = errors.New("user not authenticated")
	ErrUserStatusNotFound   = errors.New("user status not found")
	ErrEmailAlreadyExist    = errors.New("email already exist")
	ErrDuplicateUsername    = errors.New("username already exists")
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrInactiveAccount      = errors.New("account is inactive")
	ErrInvalidToken         = errors.New("invalid token")
	ErrActivationFailed     = errors.New("activation failed")
	ErrVerificationFailed   = errors.New("verification failed")
	ErrEmailSendFailed      = errors.New("failed to send email")

	// Token
	ErrRefreshTokenNotFound   = errors.New("refresh token not found")
	ErrInvalidActivationToken = errors.New("invalid or expired activation token")

	// server errors
	ErrInternalServer = errors.New("internal server error")
)
