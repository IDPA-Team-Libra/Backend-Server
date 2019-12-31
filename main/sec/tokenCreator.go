package sec

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenCreator struct {
	Username string
	Secret   []byte
}

type Response struct {
	Message        string `json:"response"`
	TokenName      string `json:"tokenName"`
	Token          string `json:"token"`
	ExpirationTime int64  `json:"expires"`
	UserData       string `json:"user"`
}

func (creator *TokenCreator) CreateToken() Response {
	expirationTime := time.Now().Add(300 * time.Minute)
	claims := &Claims{
		Username: creator.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(creator.Secret)
	if err != nil {
		return Response{}
	}
	response := Response{
		TokenName:      "auth_token",
		Token:          tokenString,
		ExpirationTime: expirationTime.Unix(),
	}
	return response
}
