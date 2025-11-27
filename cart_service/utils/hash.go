package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic("Could not hash password")
	}
	return string(bcryptedPassword), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
