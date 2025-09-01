package helpers

import (
	"errors"
	"net/mail"
	"regexp"
	"time"

	"github.com/dlclark/regexp2"
)

var (
	emailRegex    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	passwordRegex = regexp2.MustCompile(
		`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[ !"#$%&'()*+,\-./:;<=>?@[\\\]^_{|}~]).{8,}$`,
		0,
	)
)

// Validate login input
func ValidateLoginInput(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

// Validate booking creation
func ValidateCreateBookingInput(guideID uint, datetime, notes string) error {
	if guideID == 0 {
		return errors.New("guide ID is required")
	}
	if _, err := time.Parse(time.RFC3339, datetime); err != nil {
		return errors.New("invalid datetime format; must be RFC3339 (e.g., 2006-01-02T15:04:05Z07:00)")
	}
	return nil
}

// Validate registration
func ValidateRegisterInput(name, email, password string) error {
	if l := len(name); l < 3 || l > 50 {
		return errors.New("name must be between 3 and 50 characters")
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	ok, _ := passwordRegex.MatchString(password) // regexp2 returns (bool, error)
	if !ok {
		return errors.New("password must be at least 8 characters and include at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}
	return nil
}

// Email validation helper
func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email format")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}
