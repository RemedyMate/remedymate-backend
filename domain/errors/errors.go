package errors

import "errors"

// sentinel errors for content/topic mapping
var (
	ErrTopicNotFound        = errors.New("topic not found")
	ErrLanguageNotAvailable = errors.New("language not available")
	ErrUnsupportedLanguage  = errors.New("unsupported language")
	ErrNoTopicMapped        = errors.New("no topic could be mapped from the provided symptoms")
)
