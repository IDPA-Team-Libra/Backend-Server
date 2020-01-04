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
	"github.com/robfig/cron"
)

var dbConn *sql.DB

//SetupLogger creates a new logger instance
func SetupLogger() {
	logFilePath, _ := filepath.Abs(fmt.Sprintf("log/libra_log"))
	logger.SetupLogger(logFilePath)
}

//GetDatabaseInstance a "global" method to get a reference to the sql.DB
func GetDatabaseInstance() *sql.DB {
	return dbConn
}

func main() {
	// setup service with a http server
	service := service.NewService("#001", "login", "A login service that handles login for users", "3440")
	service.DefaultRoutes = false
	SetupLogger()
	defer logger.SyncLogger()
	logger.LogMessage("Server has started on 3440", logger.INFO)
	service.ActivateHTTPServer()
	//service.ActivateHTTPSServer("certs/server.crt", "certs/server.key")
	service.SetDatabaseInformation("localhost", "3306", "mysql", "administrator", "LOCAL1234", "libra")
	//service.SetDatabaseInformation("localhost", "3306", "mysql", "root", "pw123", "libra")
	db := service.GetDatabaseConnection()
	dbConn = db
	dbConn = db
	setDatabaseReferences(db)
	defer db.Close()
	mail.SetMailConfiguration(mail.Configuration{
		Sender: "librastockcompany@gmail.com",
		Pass:   "0508f610ab733eb6bc27d06587854697",
		UserID: "0d3a0e52dc9907a6ba30762a2a0999c3",
	})
	/*
		SPACE FOR ROUTES
	*/
	service.AddHTTPRoute("/user/login", Login)
	service.AddHTTPRoute("/user/register", Register)
	service.AddHTTPRoute("/user/changePassword", ChangePassword)
	service.AddHTTPRoute("/stock/all", GetStocks)
	service.AddHTTPRoute("/stock/{symbol}/{interval}", GetStockByParameter)
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

//SetupCronJobs -- creates a cronjob that executes 3 different functions that are run on different times and handle stock-fetching, cache clear and Delayed Transaction exectuion
func SetupCronJobs() {
	cronJob := cron.New()
	_, err := cronJob.AddFunc("@every 40m", func() {
		apiconnection.LoadAllStocks("5")
	})
	if err != nil {
		panic(err.Error())
	}
	_, err = cronJob.AddFunc("@every 25m", func() {
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
