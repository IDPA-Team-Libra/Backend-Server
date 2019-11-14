package main

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/Sanctuary/service"
	"github.com/Liberatys/libra-back/main/mail"
)

var database *sql.DB

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	service.ActivateHTTPServer()
	fmt.Println(service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "1234", "libra"))
	database = service.GetDatabaseConnection()
	//apiconnection.GetStockDataForSymbol("TSLA", av.TimeIntervalFiveMinute)
	mailer = mail.NewMail("mountainviewcasino@gmail.com", "1234", "Wir heissen Sie herzlich bei Libra wilkommen", "Welcome to libra", "")
	//stockapi.SendRequest()
	/*
		SPACE FOR ROUTES
	*/
	service.AddHTTPRoute("/user/login", Login)
	service.AddHTTPRoute("/user/register", Register)
	/*
		END SPACE FOR ROUTES
	*/
	service.StartHTTPServer()
}
