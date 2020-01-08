package sec

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//TokenCreator struct to hold and operate with username and secret
type TokenCreator struct {
	Username string
	Secret   []byte
}

//Response a response for the token validation or creation
type Response struct {
	Message        string `json:"response"`
	TokenName      string `json:"tokenName"`
	Token          string `json:"token"`
	ExpirationTime int64  `json:"expires"`
	UserData       string `json:"user"`
}

//NewTokenCreator creates a new TokenCreator based on secret and username
func NewTokenCreator(username string, secret []byte) TokenCreator {
	return TokenCreator{
		Username: username,
		Secret:   secret,
	}
}

const (
	//EXPIRATION marks the expiration in minutes for a given token
	EXPIRATION = 300
	//TOKENNAME the name with which the token is stored in the browser
	TOKENNAME = "auth_token"
)

//CreateToken creates a token with an expiration of [EXPIRATION] and returns it wrapped into a [Response]
func (creator *TokenCreator) CreateToken() Response {
	expirationTime := time.Now().Add(EXPIRATION * time.Minute)
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
		TokenName:      TOKENNAME,
		Token:          tokenString,
		ExpirationTime: expirationTime.Unix(),
	}
	return response
}
