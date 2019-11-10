package main

import (
	"fmt"

	"github.com/Liberatys/Sanctuary/service"
)

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.ActivateHTTPServer()
	service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "PLACEHOLDER", "libra")
	result, err := service.ExecutePerparedQuery("INSERT INTO User values(?)", "Joseph")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(service.SerializeQueryResult(result))
	/*
		SPACE FOR ROUTES
	*/

	/*
		END SPACE FOR ROUTES
	*/
	//service.StartHTTPServer()
}
