package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Liberatys/libra-back/main/logger"
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
	TotalStocks  string `json:"totalStocks"`
	StartCapital string `json:"startCapital"`
}

type Author struct {
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
	user_instance.ID = user.GetUserIdByUsername(user_instance.Username, GetDatabaseInstance())
	success, message := user_instance.Authenticate(GetDatabaseInstance())
	response := sec.Response{}
	if success == true {
		response = GenerateTokenForUser(currentUser.Username)
		logger.LogMessage(fmt.Sprintf("Nutzer hat sich eingeloggt | User %s", currentUser.Username), logger.WARNING)
	} else {
		response.Message = message
		resp, _ := json.Marshal(response)
		logger.LogMessage(message, logger.WARNING)
		w.Write(resp)
		return
	}
	response.Message = message
	currentUser.Password = ""
	portfolio_inst := user.LoadPortfolio(currentUser.Username, GetDatabaseInstance())
	currentUser.Portfolio = SerializedPortfolio{
		CurrentValue: portfolio_inst.CurrentValue.String(),
		TotalStocks:  strconv.FormatInt(portfolio_inst.TotalStocks, 10),
		StartCapital: portfolio_inst.StartCapital.String(),
	}
	user_data, _ := json.Marshal(currentUser)
	response.UserData = string(user_data)
	resp, err := json.Marshal(response)
	if err != nil {
		logger.LogMessage(message, logger.WARNING)
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
	uniqueUsername := user_instance.IsUniqueUsername(GetDatabaseInstance())
	if uniqueUsername == true {
		success, error_message := user_instance.CreationSetup(GetDatabaseInstance())
		if success == false {
			logger.LogMessage(error_message, logger.WARNING)
			w.Write([]byte(error_message))
			return
		}
		if user_instance.Write(GetDatabaseInstance()) == false {
			logger.LogMessage(fmt.Sprintf("Ungültige Daten in der Nutzer erstellung | Daten %s|%s", currentUser.Username, currentUser.Email), logger.WARNING)
			w.Write([]byte("Benutzer konnte nicht erstellt werden. Bitte an Kundendienst wenden"))
			return
		} else {
			response := sec.Response{}
			if success == true {
				response = GenerateTokenForUser(currentUser.Username)
			}
			response.Message = "Success"
			user_id := user.GetUserIdByUsername(user_instance.Username, GetDatabaseInstance())
			portfolio := user.Portfolio{}
			var accountStartBalance float64
			// if no value is set for the start balance, just take 100000 as a fall backnumber
			if currentUser.StartBalance == "" {
				accountStartBalance = START_CAPITAL
			} else {
				accountStartBalance, _ = strconv.ParseFloat(currentUser.StartBalance, 64)
			}
			portfolio.Write(user_id, GetDatabaseInstance(), accountStartBalance)
			currentUser.Password = ""
			currentUser.Portfolio = SerializedPortfolio{
				CurrentValue: portfolio_inst.CurrentValue.String(),
				TotalStocks:  strconv.FormatInt(portfolio_inst.TotalStocks, 10),
				StartCapital: portfolio_inst.StartCapital.String(),
				CurrentValue: portfolio.CurrentValue.String(),
				StartCapital: portfolio.StartCapital.String(),
			}
			user_data, _ := json.Marshal(currentUser)
			response.UserData = string(user_data)
			resp, err := json.Marshal(response)
			if err != nil {
				logger.LogMessage(fmt.Sprintf("Invalid Authentication | User %s", currentUser.Username), logger.WARNING)
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
			logger.LogMessage(fmt.Sprintf("Invalid Authentication | User %s", currentUser.Username), logger.WARNING)
		}
	} else {
		response.Message = "Valid Token"
		if err != nil {
			logger.LogMessage(fmt.Sprintf("Invalid Authentication | User %s", currentUser.Username), logger.WARNING)
		}
	}
	resp, err := json.Marshal(response)
	w.Write(resp)
}

func Logout(w http.ResponseWriter, r *http.Request) {
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
	response := sec.Response{}
	if validator.IsValidToken(jwtKey) == false {
		response.Message = "Invalid Token"
		logger.LogMessage(fmt.Sprintf("Invalid Authentication | User %s", currentUser.Username), logger.WARNING)
		if err != nil {
			logger.LogMessage(err.Error(), logger.WARNING)
		}
	} else {
		response.Message = "Valid Token"
		if err != nil {
			logger.LogMessage(err.Error(), logger.WARNING)
		}
	}
}

type PasswordChangeRequest struct {
	Username    string `json:"username"`
	AuthToken   string `json:"authToken"`
	NewPassword string `json:"newPassword"`
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var request PasswordChangeRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
		w.Write([]byte("Invalid request format"))
		return
	}
	currentUser := user.User{
		Username: request.Username,
		Password: request.NewPassword,
	}
	validator := sec.NewValidator(request.AuthToken, request.Username)
	if validator.IsValidToken(jwtKey) == false {
		obj, _ := json.Marshal("Not able to authenticate user")
		w.Write([]byte(obj))
		return
	}
	success := changePassword(currentUser.Username, currentUser.Password)
	if success == false {
		obj, _ := json.Marshal("Passwort konnte nicht geändert werden")
		w.Write([]byte(obj))
	} else {
		obj, _ := json.Marshal("Das Passwort wurde geändert")
		w.Write([]byte(obj))
	}
	return
}

func changePassword(username string, newPassword string) bool {
	password_validator := user.NewPasswordValidator(newPassword)
	pw_hash := password_validator.HashPassword()
	userID := user.GetUserIdByUsername(username, GetDatabaseInstance())
	success := user.OverwritePasswordForUserId(userID, pw_hash, database)
	return success
}
