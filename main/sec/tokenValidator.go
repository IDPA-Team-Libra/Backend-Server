package sec

import (
	"github.com/dgrijalva/jwt-go"
)

type Validator struct {
	tokenString string
	username    string
	secret_key  string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewValidator(tokenString, username string) Validator {
	validator := Validator{
		tokenString: tokenString,
		username:    username,
	}
	return validator
}

func (validator *Validator) IsValidToken(secret []byte) bool {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(validator.tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false
		}
		return false
	}
	if !tkn.Valid {
		return false
	}
	if claims.Username == validator.username {
		return true
	}
	return false
}
