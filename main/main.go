package main

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/Sanctuary/service"
	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/stock"
)

var database *sql.DB

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	service.ActivateHTTPServer()
	fmt.Println(service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "libra"))
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
	/*
		END SPACE FOR ROUTES
	*/
	//apiconnection.LoadAllStocks("5")
	service.StartHTTPServer()
}

func setDatabaseReferences(database *sql.DB) {
	stock.SetDatabaseConnection(database)
}
