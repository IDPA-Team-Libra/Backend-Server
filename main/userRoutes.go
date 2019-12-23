package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/sec"
	"github.com/Liberatys/libra-back/main/user"
)

var jwtKey = []byte("Secret")
var mailer mail.Mail

type User struct {
	Username     string              `json:"username"`
	Password     string              `json:"password"`
	Email        string              `json:"email"`
	StartBalance string              `json:"startBalance"`
	AccessToken  string              `json:"accessToken"`
	Portfolio    SerializedPortfolio `json:"portfolio"`
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
	user_instance.ID = user_instance.GetUserIdByUsername(user_instance.Username)
	success, message := user_instance.Authenticate()
	response := sec.Response{}
	if success == true {
		response = GenerateTokenForUser(currentUser.Username)
	} else {
		response.Message = message
		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
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

const (
	START_CAPITAL = 1000000
)

func GenerateTokenForUser(username string) sec.Response {
	creator := sec.TokenCreator{Username: username, Secret: jwtKey}
	return creator.CreateToken()
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
	uniqueUsername := user_instance.IsUniqueUsername()
	if uniqueUsername == true {
		success, error_message := user_instance.CreationSetup()
		if success == false {
			fmt.Println(error_message)
			return
		}
		if user_instance.Write() == false {
			w.Write([]byte("Benutzer konnte nicht erstellt werden. Bitte an Kundendienst wenden"))
		} else {
			response := sec.Response{}
			if success == true {
				response = GenerateTokenForUser(currentUser.Username)
			}
			response.Message = "Success"
			user_id := user_instance.GetUserIdByUsername(user_instance.Username)
			portfolio := user.Portfolio{}
			var accountStartBalance float64
			// if no value is set for the start balance, just take 100000 as a fall backnumber
			if currentUser.StartBalance == "" {
				accountStartBalance = 100000.0
			} else {
				accountStartBalance, _ = strconv.ParseFloat(currentUser.StartBalance, 64)
			}
			portfolio.Write(user_id, user_instance, accountStartBalance)
			currentUser.Password = ""
			currentUser.Portfolio = SerializedPortfolio{
				CurrentValue: portfolio.CurrentValue.String(),
				StartCapital: portfolio.StartCapital.String(),
			}
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
		responseObject, _ := json.Marshal("Benutzername bereits vergeben")
		w.Write(responseObject)
	}
}

func ValidateUserToken(w http.ResponseWriter, r *http.Request) {
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
	validator := sec.NewValidator(currentUser.AccessToken, currentUser.Username)
	response := PortfolioContent{}
	if validator.IsValidToken(jwtKey) == false {
		response.Message = "Invalid Token"
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		response.Message = "Valid Token"
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	resp, err := json.Marshal(response)
	w.Write(resp)
}

func Logout(w http.ResponseWriter, r *http.Request) {

}
