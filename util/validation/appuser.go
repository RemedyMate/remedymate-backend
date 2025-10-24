package validation

import (
	"regexp"
	"strings"
)

// ISO 3166-1 alpha-2 country code validation
func IsISO3166Alpha2(code string) bool {
	if len(code) != 2 {
		return false
	}
	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

// ISO 639-1 language code (2 letters, lowercased)
func IsISO6391(code string) bool {
	if len(code) != 2 {
		return false
	}
	for _, r := range code {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

// E.164 phone number: + followed by 8-15 digits
var e164 = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

func IsE164Phone(phone string) bool {
	phone = strings.TrimSpace(phone)
	return e164.MatchString(phone)
}
