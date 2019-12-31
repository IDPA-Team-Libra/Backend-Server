package main

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/Liberatys/Sanctuary/service"
	"github.com/Liberatys/libra-back/main/apiconnection"
	"github.com/Liberatys/libra-back/main/database"
	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/user"
	"github.com/robfig/cron"
)

var dbConn *sql.DB

const (
	ExMode = "TEST"
)

func SetupLogger() {
	logFilePath, _ := filepath.Abs(fmt.Sprintf("log/libra_log"))
	logger.SetupLogger(logFilePath, 4, 5)
}

func GetDatabaseInstance() *sql.DB {
	return dbConn
}

func MailMessage() string {
	return "Wir heissen Sie herzlich bei Libra wilkommen. Libra ist eine Simulierung des wirklichen Aktienmarkets und soll Ihnen helfen Aktien zu kaufen und verkaufen"
}

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	SetupLogger()
	defer logger.SyncLogger()
	logger.LogMessage("Server has started on 3440", logger.INFO)
	service.ActivateHTTPServer()
	service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "Siera_001_DB", "libra")
	//service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "pw123", "libra")
	db := service.GetDatabaseConnection()
	dbConn = db
	// db.SetMaxIdleConns(0)
	// db.SetMaxOpenConns(500)
	// db.SetConnMaxLifetime(time.Second * 10)
	dbConn = db
	setDatabaseReferences(db)
	defer db.Close()
	mailer = mail.NewMail("mountainviewcasino@gmail.com", "1234", "Wir heissen Sie herzlich bei Libra wilkommen", "Welcome to libra")
	if ExMode == "DEV" {
		userInstance := user.CreateUserInstance("Haspi", "1234", " ")
		userInstance.CreationSetup(GetDatabaseInstance())
		userInstance.Write(GetDatabaseInstance())
		userID := user.GetUserIdByUsername(userInstance.Username, GetDatabaseInstance())
		portfolio := user.Portfolio{}
		if portfolio.Write(userID, GetDatabaseInstance(), DefaultStartCapital) == false {
			logger.LogMessage("Was not able to create default user", logger.WARNING)
		} else {
		}
	}
	/*
		SPACE FOR ROUTES
	*/
	service.AddHTTPRoute("/user/login", Login)
	service.AddHTTPRoute("/user/register", Register)
	service.AddHTTPRoute("/user/changePassword", ChangePassword)
	service.AddHTTPRoute("/stock/all", GetStocks)
	service.AddHTTPRoute("/transaction/buy", AddTransaction)
	service.AddHTTPRoute("/transaction/get/delayed", GetDelayedTransactionsByUser)
	service.AddHTTPRoute("/transaction/sell", RemoveTransaction)
	service.AddHTTPRoute("/transaction/buy/delayed", AddDelayedBuyTransaction)
	service.AddHTTPRoute("/transaction/sell/delayed", AddDelayedSellTransaction)
	service.AddHTTPRoute("/transaction/all", GetUserTransaction)
	service.AddHTTPRoute("/portfolio/get", GetPortfolio)
	service.AddHTTPRoute("/authenticate/token", ValidateUserToken)
	/*
		END SPACE FOR ROUTES
	*/
	go apiconnection.LoadAllStocks("5")
	SetupCronJobs()
	service.StartHTTPServer()
}

func setDatabaseReferences(database *sql.DB) {
	stock.SetDatabaseConnection(GetDatabaseInstance())
}

func SetupCronJobs() {
	cronJob := cron.New()
	_, err := cronJob.AddFunc("@every 40m", func() {
		apiconnection.LoadAllStocks("5")
	})
	if err != nil {
		panic(err.Error())
	}
	_, err = cronJob.AddFunc("@every 20m", func() {
		PurgeStockScreen()
	})
	if err != nil {
		panic(err.Error())
	}
	_, err = cronJob.AddFunc("@midnight", func() {
		database.StartBatchProcess(GetDatabaseInstance())
	})
	if err != nil {
		panic(err.Error())
	}
	cronJob.Start()
}
