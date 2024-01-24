package util

import "golang.org/x/crypto/bcrypt"

// HashPassword returns the bcrypt hash of password
func HashPassword(password string) (string, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hasedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password, hasedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))
}
