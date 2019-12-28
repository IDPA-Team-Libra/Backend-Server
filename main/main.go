package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Liberatys/Sanctuary/service"
	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/user"
)

var database *sql.DB

const (
	EX_MODE = "TEST"
)

func setupDB() {

}

func SetupLogger() {
	now := time.Now()
	valu := now.Format("2006_01_02_15_04_05")
	log_file_path, _ := filepath.Abs(fmt.Sprintf("log/%s.txt", valu))
	logger.SetupLogger(log_file_path, 4, 5)
}

func GetDatabaseInstance() *sql.DB {
	return database
}

func MailMessage() string {
	return "Wir heissen Sie herzlich bei Libra wilkommen. Libra ist eine Simulierung des wirklichen Aktienmarkets und soll Ihnen helfen Aktien zu kaufen und verkaufen"
}

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	SetupLogger()
	logger.LogMessage("Server has started on 3440", logger.INFO)
	service.ActivateHTTPServer()
	//service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "Siera_001_DB", "libra")
	service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "pw123", "libra")
	database = service.GetDatabaseConnection()
	setDatabaseReferences(database)
	defer database.Close()
	mailer = mail.NewMail("mountainviewcasino@gmail.com", "1234", "Wir heissen Sie herzlich bei Libra wilkommen", "Welcome to libra")
	if EX_MODE == "DEV" {
		user_instance := user.CreateUserInstance("Haspi", "1234", " ")
		user_instance.CreationSetup(GetDatabaseInstance())
		user_instance.Write(GetDatabaseInstance())
		user_id := user.GetUserIdByUsername(user_instance.Username, GetDatabaseInstance())
		portfolio := user.Portfolio{}
		portfolio.Write(user_id, GetDatabaseInstance(), START_CAPITAL)
	}
	/*
		SPACE FOR ROUTES
	*/
	service.AddHTTPRoute("/user/login", Login)
	service.AddHTTPRoute("/user/register", Register)
	service.AddHTTPRoute("/user/changePassword", ChangePassword)
	service.AddHTTPRoute("/stock/all", GetStocks)
	service.AddHTTPRoute("/transaction/buy", AddTransaction)
	service.AddHTTPRoute("/transaction/sell", RemoveTransaction)
	service.AddHTTPRoute("/transaction/buy/delayed", AddDelayedTransaction)
	service.AddHTTPRoute("/transaction/all", GetUserTransaction)
	service.AddHTTPRoute("/portfolio/get", GetPortfolio)
	service.AddHTTPRoute("/authenticate/token", ValidateUserToken)
	/*
		END SPACE FOR ROUTES
	*/
	//go apiconnection.LoadAllStocks("5")
	service.StartHTTPServer()
}

func setDatabaseReferences(database *sql.DB) {
	stock.SetDatabaseConnection(database)
}
