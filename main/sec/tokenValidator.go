package sec

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

//TokenValidator holds information about the receaved token and the Username
type TokenValidator struct {
	TokenString string
	Username    string
}

//Claims claims to parse for the jwt-library
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//NewTokenValidator creates a new Tokenvalidator
func NewTokenValidator(TokenString, Username string) TokenValidator {
	TokenValidator := TokenValidator{
		TokenString: TokenString,
		Username:    Username,
	}
	return TokenValidator
}

//IsValidToken checks weather at token matches the given secret and contains the username that is given
func (TokenValidator *TokenValidator) IsValidToken(secret []byte) bool {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(TokenValidator.TokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false
		}
		fmt.Println(err)
		return false
	}
	if !tkn.Valid {
		return false
	}
	if claims.Username == TokenValidator.Username {
		return true
	}
	return false
}
