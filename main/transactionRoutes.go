package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/Liberatys/libra-back/main/database"
	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/sec"
	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

type TransactionRequest struct {
	AuthToken          string `json:"authToken"`
	Username           string `json:"username"`
	StockSymbol        string `json:"stockSymbol"`
	Operation          string `json:"operation"`
	Amount             int64  `json:"amount"`
	Date               string `json:"date"`
	ExpectedStockPrice string `json:"expectedStockPrice"`
}

type TransactionResponse struct {
	Message   string `json:"message"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Operation string `json:"operation"`
	Value     string `json:"transactionValue"`
}

type FutureTransactionOption struct {
	TransactionRequest TransactionRequest `json:"transactionRequest"`
	SetDate            string             `json:"setDate"`
}

func GetUserTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var request TransactionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.Write([]byte("Invalid request format"))
		return
	}
	currentUser := user.User{
		Username: request.Username,
	}
	validator := sec.NewValidator(request.AuthToken, currentUser.Username)
	if validator.IsValidToken(jwtKey) == false {
		logger.LogMessage(fmt.Sprintf("Anfrage an GetUserTransaction hatte einen ungültigen jwt. | User: %s", currentUser.Username), logger.WARNING)
		response := PortfolioContent{}
		response.Message = "Invalid Token"
		resp, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Write(resp)
		return
	}
	//TODO do something
	if err != nil {
		fmt.Println(err.Error())
	}
	userID := user.GetUserIdByUsername(request.Username, GetDatabaseInstance())
	trans := transaction.Transaction{}
	transactions := trans.LoadTransactionsByProcessState(userID, GetDatabaseInstance(), true)
	json_obj, err := json.Marshal(transactions)
	if err != nil {
		fmt.Println(err.Error())
		logger.LogMessage(fmt.Sprintf("Das Request Format in einer Anfrage an GetUserTransaction wurde nicht eingehalten | User: %s", currentUser.Username), logger.WARNING)
		obj, _ := json.Marshal("Invalid request format")
		w.Write(obj)
		return
	}
	w.Write(json_obj)
}

func GetDelayedTransactionsByUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var request TransactionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.Write([]byte("Invalid request format"))
		return
	}
	currentUser := user.User{
		Username: request.Username,
	}
	validator := sec.NewValidator(request.AuthToken, currentUser.Username)
	if validator.IsValidToken(jwtKey) == false {
		logger.LogMessage(fmt.Sprintf("Anfrage an GetUserTransaction hatte einen ungültigen jwt. | User: %s", currentUser.Username), logger.WARNING)
		response := PortfolioContent{}
		response.Message = "Invalid Token"
		resp, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Write(resp)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	userID := user.GetUserIdByUsername(request.Username, GetDatabaseInstance())
	trans := transaction.Transaction{}
	transactions := trans.LoadTransactionsByProcessState(userID, GetDatabaseInstance(), false)
	json_obj, err := json.Marshal(transactions)
	if err != nil {
		fmt.Println(err.Error())
		logger.LogMessage(fmt.Sprintf("Das Request Format in einer Anfrage an GetUserTransaction wurde nicht eingehalten | User: %s", currentUser.Username), logger.WARNING)
		obj, _ := json.Marshal("Invalid request format")
		w.Write(obj)
		return
	}
	w.Write(json_obj)
}

func RemoveTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		obj, _ := json.Marshal("Keine Parameter übergeben")
		w.Write(obj)
		return
	}
	var request TransactionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		obj, _ := json.Marshal("Invalid request format")
		w.Write(obj)
		return
	}
	currentUser := user.User{
		Username: request.Username,
	}
	validator := sec.NewValidator(request.AuthToken, request.Username)
	if validator.IsValidToken(jwtKey) == false {
		logger.LogMessageWithOrigin(fmt.Sprintf("User: %s was not able to be authenticated", currentUser.Username), logger.WARNING, "transactionRoutes.go")
		response := TransactionResponse{
			Message:   "Leider konnten Sie nicht durch den Server authentizifiert werden. Bitte neu einloggen",
			State:     "Breach",
			Title:     "Aktion konte nicht ausgeführt werden",
			Operation: "-",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	user_instance := user.CreateUserInstance(currentUser.Username, currentUser.Password, "")
	user_instance.ID = user.GetUserIdByUsername(request.Username, GetDatabaseInstance())
	requestedStock := stock.LoadStockInstance(request.StockSymbol)
	items := user.LoadUserItems(user_instance.ID, request.StockSymbol, GetDatabaseInstance())
	totalStockQuantity := database.CalculateTotalStocks(items)
	requestCount := request.Amount
	if totalStockQuantity < requestCount {
		logger.LogMessageWithOrigin(fmt.Sprintf("User: %s user tried to sell more stock then he could", currentUser.Username), logger.WARNING, "transactionRoutes.go")
		response := TransactionResponse{
			Message:   "Sie können nicht mehr Aktien verkaufen als sie haben",
			State:     "Failed",
			Title:     "Verkauf konnte nicht durchgeführt werden",
			Operation: "-",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	handler, err := GetDatabaseInstance().Begin()
	items = database.SubtractStocksFromTotalAmount(items, requestCount, handler)
	//TODO do something
	if err != nil {
		obj, _ := json.Marshal("Invalid request format")
		w.Write(obj)
		return
	}
	transaction := transaction.NewTransaction(user_instance.ID, request.Operation, request.Operation+" "+request.StockSymbol, request.Amount, requestedStock.Price, request.Date)
	if transaction.Write(true, handler) == false {
		handler.Rollback()
		obj, _ := json.Marshal("Invalid request format")
		w.Write(obj)
		return
	}
	portfolio := user.LoadPortfolio(request.Username, GetDatabaseInstance())
	portfolio.TotalStocks -= request.Amount
	s := fmt.Sprintf("%f", float64(request.Amount))
	additionalBalance := database.MultiplyString(s, requestedStock.Price)
	portfolio.Balance = *portfolio.Balance.Add(&portfolio.Balance, additionalBalance)
	portfolio.CurrentValue = *portfolio.CurrentValue.Sub(&portfolio.CurrentValue, additionalBalance)
	if portfolio.Update(handler) == false {
		handler.Rollback()
		response := TransactionResponse{
			Message:   "Verkauf konnte nicht getätigt werden",
			State:     "Failure",
			Title:     "Kauf bitte erneut ausführen",
			Operation: "-",
			Value:     "0",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	handler.Commit()
	logger.LogMessageWithOrigin(fmt.Sprintf("User: %s user sold %s", currentUser.Username, request.StockSymbol), logger.WARNING, "transactionRoutes.go")
	response := TransactionResponse{
		Message:   "Verkauf wurde getätigt",
		State:     "Success",
		Title:     "Titel wurde verkauft und der Betrag ihrem Konto gutgeschrieben",
		Operation: "-",
		Value:     additionalBalance.String(),
	}
	obj, _ := json.Marshal(response)
	w.Write([]byte(obj))
	return
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var request TransactionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("Invalid request format"))
		return
	}
	currentUser := user.User{
		Username: request.Username,
	}
	userID := user.GetUserIdByUsername(request.Username, GetDatabaseInstance())
	currentUser.ID = userID
	if userID <= 0 {
		fmt.Println("Invalud userID")
		return
	}
	validator := sec.NewValidator(request.AuthToken, request.Username)
	if validator.IsValidToken(jwtKey) == false {
		logger.LogMessageWithOrigin(fmt.Sprintf("User: %s was not able to be authenticated", currentUser.Username), logger.WARNING, "transactionRoutes.go | AddTransaction")
		response := TransactionResponse{
			Message:   "Leider konnten Sie nicht durch den Server authentizifiert werden. Bitte neu einloggen",
			State:     "Breach",
			Title:     "Aktion konte nicht ausgeführt werden",
			Operation: "+",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	portfolio := user.LoadPortfolio(request.Username, GetDatabaseInstance())
	totalPrice := new(big.Float)
	requestedStock := stock.LoadStockInstance(request.StockSymbol)
	totalPrice, _ = totalPrice.SetString(requestedStock.Price)
	amount := float64(request.Amount)
	totalPrice = totalPrice.Mul(totalPrice, big.NewFloat(amount))
	if portfolio.Balance.Cmp(totalPrice) != 1 {
		logger.LogMessageWithOrigin(fmt.Sprintf("User: %s tried to sell more than he had", currentUser.Username), logger.WARNING, "transactionRoutes.go | AddTransaction")
		response := TransactionResponse{
			Message:   "Dieser Kauf überschreitet leider Ihren Kontostand",
			State:     "Failed",
			Title:     "Kauf konnte nicht durchgeführt werden",
			Operation: "+",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	transaction := transaction.NewTransaction(userID, request.Operation, request.Operation+" "+request.StockSymbol, request.Amount, requestedStock.Price, request.Date)
	database.CreateTransaction(transaction, portfolio, requestedStock, currentUser, request.Amount, *totalPrice, GetDatabaseInstance())
	logger.LogMessageWithOrigin(fmt.Sprintf("User: %s bought %s", currentUser.Username, request.StockSymbol), logger.WARNING, "transactionRoutes.go | AddTransaction")
	response := TransactionResponse{
		Message:   "Kauf wird abgewickelt.. Dies kann je nach Auslastung einige Minuten dauern",
		State:     "Success",
		Title:     "Kauf abgeschlossen",
		Operation: "+",
		Value:     totalPrice.String(),
	}
	obj, _ := json.Marshal(response)
	w.Write([]byte(obj))
}

func HandleDelayedTransaction(w http.ResponseWriter, r *http.Request) (bool, transaction.Transaction) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return false, transaction.Transaction{}
	}
	var request TransactionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
		w.Write([]byte("Invalid request format"))

		return false, transaction.Transaction{}
	}
	currentUser := user.User{
		Username: request.Username,
	}
	userID := user.GetUserIdByUsername(request.Username, GetDatabaseInstance())
	currentUser.ID = userID
	if userID <= 0 {
		logger.LogMessage(fmt.Sprintf("Eine Nutzer ID war nicht gültig | User %s", request.Username), logger.WARNING)
		return false, transaction.Transaction{}
	}
	validator := sec.NewValidator(request.AuthToken, request.Username)
	if validator.IsValidToken(jwtKey) == false {
		logger.LogMessageWithOrigin(fmt.Sprintf("User: %s was not able to be validated", currentUser.Username), logger.WARNING, "transactionRoutes.go | AddDelayedTransaction")
		response := TransactionResponse{
			Message:   "Leider konnten Sie nicht durch den Server authentizifiert werden. Bitte neu einloggen",
			State:     "Breach",
			Title:     "Aktion konte nicht ausgeführt werden",
			Operation: "-",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return false, transaction.Transaction{}
	}
	totalPrice := new(big.Float)
	requestedStock := stock.LoadStockInstance(request.StockSymbol)
	totalPrice, _ = totalPrice.SetString(requestedStock.Price)
	amount := float64(request.Amount)
	totalPrice = totalPrice.Mul(totalPrice, big.NewFloat(amount))
	transaction := transaction.NewTransaction(userID, request.Operation, request.Operation+" "+request.StockSymbol, request.Amount, request.ExpectedStockPrice, request.Date)
	return true, transaction
}

func AddDelayedBuyTransaction(w http.ResponseWriter, r *http.Request) {
	authenticated, transaction := HandleDelayedTransaction(w, r)
	if authenticated == false {
		return
	}
	handler, err := GetDatabaseInstance().Begin()
	if err != nil {

	}
	if transaction.Write(false, handler) == false {
		response := TransactionResponse{
			Message:   "Kauf konnte nicht abgewickelt werden",
			State:     "Failure",
			Title:     "Kauf abgeschlossen.",
			Operation: "+",
		}
		handler.Rollback()
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))

	}
	response := TransactionResponse{
		Message:   "Kaufaktion wird eingeleitet. Kauf wird am angegeben Datum eingeleitet.",
		State:     "Success",
		Title:     "Kauf abgeschlossen.",
		Operation: "*",
	}
	handler.Commit()
	obj, _ := json.Marshal(response)
	w.Write([]byte(obj))
}

func AddDelayedSellTransaction(w http.ResponseWriter, r *http.Request) {
	authenticated, transaction := HandleDelayedTransaction(w, r)
	if authenticated == false {
		return
	}
	handler, err := GetDatabaseInstance().Begin()
	if err != nil {

	}
	if transaction.Write(false, handler) == false {
		response := TransactionResponse{
			Message:   "Verkauf konnte nicht abgewickelt werden",
			State:     "Failure",
			Title:     "Kauf abgeschlossen.",
			Operation: "-",
		}
		handler.Rollback()
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	response := TransactionResponse{
		Message:   "Verkauf wird eingeleitet. Verkauf wird am angegeben Datum eingeleitet.",
		State:     "Success",
		Title:     "Kauf abgeschlossen.",
		Operation: "*",
	}
	handler.Commit()
	obj, _ := json.Marshal(response)
	w.Write([]byte(obj))
}
