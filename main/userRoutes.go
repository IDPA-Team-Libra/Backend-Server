package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Liberatys/libra-back/main/user"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Auther struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var currentUser User
	err = json.Unmarshal(body, &currentUser)
	if err != nil {
		w.Write([]byte("Invalid json"))
		return
	}
	user := user.CreateUserInstance(currentUser.Username, currentUser.Password, "")
	user.SetDatabaseConnection(database)
	success, message := user.Authenticate()
	response := Response{}
	if success == true {
		response = GenerateTokenForUser(currentUser.Username, w)
	}
	response.Message = message
	resp, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write(resp)
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type Response struct {
	Message        string `json:"response"`
	TokenName      string `json:"tokenName"`
	Token          string `json:"token"`
	ExpirationTime int64  `json:"expires"`
}

var jwtKey = []byte("PLACEHOLDER")

func GenerateTokenForUser(username string, w http.ResponseWriter) Response {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return Response{}
	}
	response := Response{
		TokenName:      "auth_token",
		Token:          tokenString,
		ExpirationTime: expirationTime.Unix(),
	}
	return response
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var currentUser User
	err = json.Unmarshal(body, &currentUser)
	if err != nil {
		w.Write([]byte("Invalid json"))
		return
	}
	user := user.CreateUserInstance(currentUser.Username, currentUser.Password, currentUser.Email)
	user.SetDatabaseConnection(database)
	if user.IsUniqueUsername() == true {
		success, error_message := user.CreationSetup()
		if success == false {
			fmt.Println(error_message)
			return
		}
		if user.Write() == false {
			w.Write([]byte("Benutzer konnte nicht erstellt werden. Bitte an Kundendienst wenden"))
		} else {
			response := Response{}
			if success == true {
				response = GenerateTokenForUser(currentUser.Username, w)
			}
			response.Message = "1"
			resp, err := json.Marshal(response)
			if err != nil {
				fmt.Println(err.Error())
			}
			w.Write(resp)
		}
	} else {
		w.Write([]byte("Benutzername bereits vergeben"))
	}
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {

}

func Logout(w http.ResponseWriter, r *http.Request) {

}
