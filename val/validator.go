package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidateUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidateFullName = regexp.MustCompile(`^[a-zA-Z]+(?:\s[a-zA-Z]+)*$`).MatchString
)

func ValidateString(str string, min, max int) error {
	n := len(str)
	if n < min || n > max {
		return fmt.Errorf("length of string must at least %d and at most %d", min, max)
	}
	return nil
}

func ValidateUsername(username string) error {
	err := ValidateString(username, 3, 20)
	if err != nil {
		return err
	}
	if !isValidateUsername(username) {
		return fmt.Errorf("must contain only lowercase letters, numbers or underscores")
	}
	return nil
}

func ValidateFullName(fullName string) error {
	err := ValidateString(fullName, 3, 100)
	if err != nil {
		return err
	}
	if !isValidateFullName(fullName) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail(email string) error {
	err := ValidateString(email, 3, 200)
	if err != nil {
		return err
	}
	_, err = mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("is not a valid email")
	}
	return nil
}
