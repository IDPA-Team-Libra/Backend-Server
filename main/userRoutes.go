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

type User struct {
	Username     string              `json:"username"`
	Password     string              `json:"password"`
	Email        string              `json:"email"`
	StartBalance string              `json:"startBalance"`
	AccessToken  string              `json:"accessToken"`
	Portfolio    SerializedPortfolio `json:"portfolio"`
}

type SerializedPortfolio struct {
	CurrentBalance string `json:"currentBalance"`
	CurrentValue   string `json:"currentValue"`
	TotalStocks    string `json:"totalStocks"`
	StartCapital   string `json:"startCapital"`
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
	userInstance := user.CreateUserInstance(currentUser.Username, currentUser.Password, "")
	_, value := userInstance.GetPasswordHashByUsername(GetDatabaseInstance())
	userInstance.ID = user.GetUserIDByUsername(userInstance.Username, GetDatabaseInstance())
	success, message := userInstance.Authenticate(GetDatabaseInstance(), value)
	fmt.Println(success)
	response := sec.Response{}
	if success == true {
		response = GenerateTokenForUser(currentUser.Username)
		logger.LogMessage(fmt.Sprintf("Nutzer hat sich eingeloggt | User %s", currentUser.Username), logger.WARNING)
	} else {
		response.Message = message
		resp, _ := json.Marshal(response)
		logger.LogMessage(fmt.Sprintf("%s | User: %s", message, currentUser.Username), logger.WARNING)
		w.Write(resp)
		return
	}
	response.Message = message
	currentUser.Password = ""
	portfolioInstance := user.LoadPortfolio(userInstance.ID, GetDatabaseInstance())
	currentUser.Portfolio = ConvertPortfolioToSerialized(portfolioInstance)
	userData, _ := json.Marshal(currentUser)
	response.UserData = string(userData)
	resp, err := json.Marshal(response)
	if err != nil {
		logger.LogMessage(message, logger.WARNING)
	}
	w.Write(resp)
}

func ConvertPortfolioToSerialized(portfolioInstance user.Portfolio) SerializedPortfolio {
	serializedPortfolio := SerializedPortfolio{
		CurrentValue:   portfolioInstance.CurrentValue.String(),
		StartCapital:   portfolioInstance.StartCapital.String(),
		TotalStocks:    strconv.FormatInt(portfolioInstance.TotalStocks, 10),
		CurrentBalance: portfolioInstance.Balance.String(),
	}
	return serializedPortfolio
}

const (
	DefaultStartCapital = 1000000
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
	userInstance := user.CreateUserInstance(currentUser.Username, currentUser.Password, currentUser.Email)
	valid, message := userInstance.Create(GetDatabaseInstance())
	if valid == true {
		response := GenerateTokenForUser(currentUser.Username)
		response.Message = "Success"
		userID := user.GetUserIDByUsername(userInstance.Username, GetDatabaseInstance())
		portfolio := user.Portfolio{}
		var accountStartBalance float64
		// if no value is set for the start balance, just take 100000 as a fall backnumber
		if currentUser.StartBalance == "" {
			accountStartBalance = DefaultStartCapital
		} else {
			accountStartBalance, _ = strconv.ParseFloat(currentUser.StartBalance, 64)
		}
		portfolio.Write(userID, GetDatabaseInstance(), accountStartBalance)
		currentUser.Password = ""
		currentUser.Portfolio = ConvertPortfolioToSerialized(portfolio)
		userData, _ := json.Marshal(currentUser)
		response.UserData = string(userData)
		resp, err := json.Marshal(response)
		if err != nil {
			logger.LogMessage(fmt.Sprintf("Invalid Authentication | User %s", currentUser.Username), logger.WARNING)
			return
		}
		mailer := mail.NewMail(currentUser.Email)
		mailer.ApplyConfiguration(mail.LoadMailConfiguration())
		go mailer.SendWelcomeEmail()
		w.Write(resp)
	} else {
		response := sec.Response{
			Message: message,
		}
		responseObject, _ := json.Marshal(response)
		w.Write(responseObject)
		return
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
	validator := sec.NewTokenValidator(currentUser.AccessToken, currentUser.Username)
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
	validator := sec.NewTokenValidator(currentUser.AccessToken, currentUser.Username)
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
	validator := sec.NewTokenValidator(request.AuthToken, request.Username)
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
		logger.LogMessage(fmt.Sprintf("Nutzer hat Passwort geändert | %s", currentUser.Username), logger.INFO)
		obj, _ := json.Marshal("Das Passwort wurde geändert")
		w.Write([]byte(obj))
	}
	return
}

func changePassword(username string, newPassword string) bool {
	passwordValidator := user.NewPasswordValidator(newPassword, "")
	passwordHash := passwordValidator.HashPassword()
	userID := user.GetUserIDByUsername(username, GetDatabaseInstance())
	success := user.OverwritePasswordForUserID(userID, passwordHash, GetDatabaseInstance())
	return success
}
