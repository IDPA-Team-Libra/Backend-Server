package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Liberatys/libra-back/main/sec"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

type PortfolioContent struct {
	Message      string `json:"response"`
	Items        string `json:"items"`
	Transactions string `json:"transactions"`
}

func GetPortfolio(w http.ResponseWriter, r *http.Request) {
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
	if validator.IsValidToken(jwtKey) == false {
		response := PortfolioContent{}
		response.Message = "Invalid Token"
		resp, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Write(resp)
		return
	}
	user_instance := user.CreateUserInstance(currentUser.Username, currentUser.Password, "")
	user_instance.SetDatabaseConnection(database)
	user_instance.ID = user_instance.GetUserIdByUsername(user_instance.Username)
	trans := transaction.Transaction{}
	trans.DatabaseConnection = database
	transactions := trans.LoadTransactions(user_instance.ID)
	response := PortfolioContent{}
	response.Message = "Success"
	item_data, _ := json.Marshal(user.LoadUserItems(user_instance, "*"))
	transaction_data, _ := json.Marshal(transactions)
	response.Items = string(item_data)
	response.Transactions = string(transaction_data)
	resp, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write(resp)
}
