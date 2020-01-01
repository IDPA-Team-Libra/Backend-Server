package user

import (
	"golang.org/x/crypto/bcrypt"
)

//PasswordValidatior structure to hold a password and to make operations on
type PasswordValidatior struct {
	Password    string
	PassworHash string
}

//NewPasswordValidator creates a new PasswordValidator to work with
func NewPasswordValidator(password, passwordHash string) PasswordValidatior {
	return PasswordValidatior{
		Password:    password,
		PassworHash: passwordHash,
	}
}

//HashPassword hashes a given password with bcrypt
func (passwordValidator *PasswordValidatior) HashPassword() string {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordValidator.Password), 11)
	if err != nil {
		return ""
	}
	return string(hash[:])
}

//ComparePasswords checks if a given password, matches a given password-hash
func (passwordValidator *PasswordValidatior) ComparePasswords() bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordValidator.PassworHash), []byte(passwordValidator.Password))
	if err != nil {
		return false
	}
	return true
}
