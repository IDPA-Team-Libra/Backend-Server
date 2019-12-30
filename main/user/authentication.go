package user

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type PasswordValidatior struct {
	Password string
}

func NewPasswordValidator(password string) PasswordValidatior {
	return PasswordValidatior{
		Password: password,
	}
}

func (passwordValidator *PasswordValidatior) isValidPassword() bool {
	if len(passwordValidator.Password) < 7 {
		return false
	}
	if strings.ContainsAny(passwordValidator.Password, "@#$%^&!.-,") == false {
		return false
	}
	if strings.ContainsAny(passwordValidator.Password, "1234567890") == false {
		return false
	}
	return true
}

func (passwordValidator *PasswordValidatior) HashPassword() string {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordValidator.Password), 11)
	if err != nil {
		//TODO: Log information to a specified logger for the package
		return ""
	}
	return string(hash[:])
}

func (passwordValidator *PasswordValidatior) comparePasswords(passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordValidator.Password))
	if err != nil {
		return false
	}
	return true
}
