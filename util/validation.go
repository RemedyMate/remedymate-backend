package util

import (
	"fmt"

	derrors "remedymate-backend/domain/AppError"
)

// ValidateLanguage checks supported language codes.
func ValidateLanguage(lang string) error {
	if lang == "" {
		return fmt.Errorf("language is required")
	}
	if lang != "en" && lang != "am" {
		return derrors.ErrUnsupportedLanguage
	}
	return nil
}

// ValidateTopicKey ensures topic key is present.
func ValidateTopicKey(topicKey string) error {
	if topicKey == "" {
		return derrors.ErrNoTopicMapped
	}
	return nil
}
