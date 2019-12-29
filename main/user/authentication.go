package user

import (
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
	//TODO: Find a good way to validate the length and composition of a passwords
	/*
		TODO: Specify the password requirements and implement the todo above.
	*/
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
