package main

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/Sanctuary/service"
)

var database *sql.DB

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	service.ActivateHTTPServer()
	fmt.Println(service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "1234", "libra"))
	database = service.GetDatabaseConnection()
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
