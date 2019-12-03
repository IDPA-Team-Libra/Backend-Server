package main

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/Sanctuary/service"
	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/stock"
)

var database *sql.DB

//TODO
//! Rewrite some of the enitites, like user
//! Write tests for the entities
//! Create a model for the backend / class diagramm
//! Refactor written code and make it more modular

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	service.ActivateHTTPServer()
	fmt.Println(service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "Siera_001_DB", "libra"))
	database = service.GetDatabaseConnection()
	setDatabaseReferences(database)
	//apiconnection.GetStockDataForSymbol("TSLA", av.TimeIntervalFiveMinute)
	mailer = mail.NewMail("mountainviewcasino@gmail.com", "1234", "Wir heissen Sie herzlich bei Libra wilkommen", "Welcome to libra")
	//stockapi.SendRequest()
	/*
		SPACE FOR ROUTES
	*/
	service.AddHTTPRoute("/user/login", Login)
	service.AddHTTPRoute("/user/register", Register)
	service.AddHTTPRoute("/stock/all", GetStocks)
	service.AddHTTPRoute("/transaction/buy", AddTransaction)
	service.AddHTTPRoute("/transaction/buy/delayed", AddDelayedTransaction)
	/*
		END SPACE FOR ROUTES
	*/
	//go apiconnection.LoadAllStocks("5")
	service.StartHTTPServer()
}

func setDatabaseReferences(database *sql.DB) {
	stock.SetDatabaseConnection(database)
}
