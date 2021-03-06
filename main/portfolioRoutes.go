package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/sec"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

//PortfolioContent holds the data sent to the client
type PortfolioContent struct {
	Message      string `json:"response"`
	Items        string `json:"items"`
	Transactions string `json:"transactions"`
}

//PortfolioRequest holds request information for the portfolio request
type PortfolioRequest struct {
	Username  string `json:"username"`
	AuthToken string `json:"authToken"`
}

//GetPortfolio get portfolio-information for a user
func GetPortfolio(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		logger.LogMessage("Anfrage an Portfolio hatte keine Parameter", logger.INFO)
		return
	}
	var currentUser PortfolioRequest
	err = json.Unmarshal(body, &currentUser)
	if err != nil {
		logger.LogMessage("Anfrage an Portfolio hatte invalides JSON", logger.INFO)
		w.Write([]byte("Invalid json"))
		return
	}
	validator := sec.NewTokenValidator(currentUser.AuthToken, currentUser.Username)
	if validator.IsValidToken(jwtKey) == false {
		logger.LogMessage(fmt.Sprintf("Anfrage an GetUserTransaction hatte einen ungültigen jwt. | User: %s", currentUser.Username), logger.WARNING)
		response := PortfolioContent{}
		response.Message = "Invalid Token"
		resp, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Write(resp)
		return
	}
	userInstance := user.CreateUserInstance(currentUser.Username, "", "")
	userInstance.ID = user.GetUserIDByUsername(userInstance.Username, GetDatabaseInstance())
	trans := transaction.Transaction{}
	transactions := trans.LoadTransactionsByProcessState(userInstance.ID, GetDatabaseInstance(), true)
	response := PortfolioContent{}
	response.Message = "Success"
	itemData, _ := json.Marshal(user.LoadUserItems(userInstance.ID, "*", GetDatabaseInstance()))
	transactionData, _ := json.Marshal(transactions)
	response.Items = string(itemData)
	response.Transactions = string(transactionData)
	resp, err := json.Marshal(response)
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
	}
	w.Write(resp)
}
