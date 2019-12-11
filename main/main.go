package main

import (
	"database/sql"

	"github.com/Liberatys/Sanctuary/service"
	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/user"
)

var database *sql.DB

//TODO
//! Rewrite some of the enitites, like user
//! Write tests for the entities
//! Create a model for the backend / class diagramm
//! Refactor written code and make it more modular

const (
	EX_MODE = "FALSE"
)

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	service.ActivateHTTPServer()
	service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "Siera_001_DB", "libra")
	database = service.GetDatabaseConnection()
	setDatabaseReferences(database)
	mailer = mail.NewMail("mountainviewcasino@gmail.com", "1234", "Wir heissen Sie herzlich bei Libra wilkommen", "Welcome to libra")
	/*
		SPACE FOR ROUTES
	*/
	if EX_MODE == "DEV" {
		user_instance := user.CreateUserInstance("Haspi", "1234", " ")
		user_instance.SetDatabaseConnection(database)
		user_instance.CreationSetup()
		user_instance.Write()
		user_id := user_instance.GetUserIdByUsername(user_instance.Username)
		portfolio := user.Portfolio{}
		portfolio.Write(user_id, user_instance, START_CAPITAL)
	}
	service.AddHTTPRoute("/user/login", Login)
	service.AddHTTPRoute("/user/register", Register)
	service.AddHTTPRoute("/stock/all", GetStocks)
	service.AddHTTPRoute("/transaction/buy", AddTransaction)
	service.AddHTTPRoute("/transaction/buy/delayed", AddDelayedTransaction)
	service.AddHTTPRoute("/transaction/all", GetUserTransaction)
	/*
		END SPACE FOR ROUTES
	*/
	//go apiconnection.LoadAllStocks("5")
	service.StartHTTPServer()
}

func setDatabaseReferences(database *sql.DB) {
	stock.SetDatabaseConnection(database)
}
