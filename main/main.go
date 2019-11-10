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
	fmt.Println(service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "PLACEHOLDER", "libra"))
	database = service.GetDatabaseConnection()
	/*
		SPACE FOR ROUTES
	*/
	service.AddHTTPRoute("/user/login", Login)
	/*
		END SPACE FOR ROUTES
	*/
	service.StartHTTPServer()
}
