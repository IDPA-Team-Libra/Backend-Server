package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/user"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("PLACEHOLDER")
var mailer mail.Mail

type User struct {
	Username  string              `json:"username"`
	Password  string              `json:"password"`
	Email     string              `json:"email"`
	Portfolio SerializedPortfolio `json:"portfolio"`
}

type SerializedPortfolio struct {
	CurrentValue string `json:"currentValue"`
	Stocks       string `json:"stocks"`
	StartCapital string `json:"startCapital"`
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
	user_instance := user.CreateUserInstance(currentUser.Username, currentUser.Password, "")
	user_instance.SetDatabaseConnection(database)
	success, message := user_instance.Authenticate()
	response := Response{}
	if success == true {
		response = GenerateTokenForUser(currentUser.Username, w)
	}
	response.Message = message
	currentUser.Password = ""
	portfolio := user.LoadPortfolio(user_instance)
	currentUser.Portfolio = SerializedPortfolio{
		CurrentValue: portfolio.CurrentValue.String(),
		StartCapital: portfolio.StartCapital.String(),
	}
	user_data, _ := json.Marshal(currentUser)
	response.UserData = string(user_data)
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
	UserData       string `json:"user"`
}

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
	user_instance := user.CreateUserInstance(currentUser.Username, currentUser.Password, currentUser.Email)
	user_instance.SetDatabaseConnection(database)
	if user_instance.IsUniqueUsername() == true {
		success, error_message := user_instance.CreationSetup()
		if success == false {
			fmt.Println(error_message)
			return
		}
		if user_instance.Write() == false {
			w.Write([]byte("Benutzer konnte nicht erstellt werden. Bitte an Kundendienst wenden"))
		} else {
			response := Response{}
			if success == true {
				response = GenerateTokenForUser(currentUser.Username, w)
			}
			response.Message = "Success"
			user_id := user_instance.GetUserIdByUsername(user_instance.Username)
			portfolio := user.Portfolio{}
			portfolio.Create(user_id, user_instance, 5000.0)
			currentUser.Password = ""
			currentUser.Portfolio = SerializedPortfolio{
				CurrentValue: portfolio.CurrentValue.String(),
				StartCapital: portfolio.StartCapital.String(),
			}
			fmt.Println(currentUser.Portfolio)
			user_data, _ := json.Marshal(currentUser)
			response.UserData = string(user_data)
			resp, err := json.Marshal(response)
			// activite if password is set and production is reached
			//go mailer.SendEmail(currentUser.Email)
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
