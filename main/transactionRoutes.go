package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/sec"
	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

//TODO implement security and user validation for this part of the system

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
	userID := user.GetUserIdByUsername(request.Username, GetDatabaseInstance())
	trans := transaction.Transaction{}
	trans.DatabaseConnection = database
	transactions := trans.LoadTransactions(userID)
	json_obj, err := json.Marshal(transactions)
	if err != nil {
		fmt.Println(err.Error())
		logger.LogMessage(fmt.Sprintf("Das Request format in einer Anfrage an GetUserTransaction wurde nicht eingehalten | User: %s", currentUser.Username), logger.WARNING)
		w.Write([]byte("Invalid request format"))
		return
	}
	w.Write(json_obj)
}

func RemoveTransaction(w http.ResponseWriter, r *http.Request) {
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
	validator := sec.NewValidator(request.AuthToken, request.Username)
	if validator.IsValidToken(jwtKey) == false {
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
	requestedStock := loadStockInstance(request.StockSymbol)
	items := user.LoadUserItems(user_instance.ID, request.StockSymbol, GetDatabaseInstance())
	totalStockQuantity := calculateTotalStocks(items)
	requestCount := request.Amount
	if totalStockQuantity < requestCount {
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
	for index, _ := range items {
		if requestCount > 0 {
			quantity := items[index].Quantity
			if quantity <= requestCount {
				requestCount -= quantity
				items[index].Quantity = 0
				items[index].Remove(GetDatabaseInstance())
			} else {
				items[index].Quantity -= requestCount
				requestCount = 0
				items[index].Update(GetDatabaseInstance())
				break
			}
		} else {
			break
		}
	}
	transaction := transaction.NewTransaction(user_instance.ID, request.Operation, request.Operation+" "+request.StockSymbol, request.Amount, requestedStock.Price, request.Date)
	transaction.DatabaseConnection = database
	transaction.Write(true)
	portfolio := user.LoadPortfolio(request.Username, GetDatabaseInstance())
	portfolio.TotalStocks -= request.Amount
	s := fmt.Sprintf("%f", float64(request.Amount))
	additionalBalance := multiplyString(s, requestedStock.Price)
	portfolio.Balance = *portfolio.Balance.Add(&portfolio.Balance, additionalBalance)
	portfolio.CurrentValue = *portfolio.CurrentValue.Sub(&portfolio.CurrentValue, additionalBalance)
	portfolio.Update(GetDatabaseInstance())
	response := TransactionResponse{
		Message:   "Verkauf wurde getätigt",
		State:     "Success",
		Title:     "Aktien wurden verkauf und ihrem Konto gutgeschrieben",
		Operation: "-",
	}
	obj, _ := json.Marshal(response)
	w.Write([]byte(obj))
	return
}

func multiplyString(first, second string) *big.Float {
	firstFloat := new(big.Float)
	firstFloat.SetString(first)
	secondFloat := new(big.Float)
	firstFloat, _ = firstFloat.SetString(first)
	secondFloat, _ = secondFloat.SetString(second)
	return firstFloat.Mul(firstFloat, secondFloat)
}

func calculateTotalStocks(items []user.PortfolioItem) int64 {
	var counter int64
	for _, stock := range items {
		counter += stock.Quantity
	}
	return counter
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
	portfolio := user.LoadPortfolio(request.Username, GetDatabaseInstance())
	totalPrice := new(big.Float)
	requestedStock := loadStockInstance(request.StockSymbol)
	totalPrice, _ = totalPrice.SetString(requestedStock.Price)
	amount := float64(request.Amount)
	totalPrice = totalPrice.Mul(totalPrice, big.NewFloat(amount))
	if portfolio.Balance.Cmp(totalPrice) != 1 {
		response := TransactionResponse{
			Message:   "Dieser Kauf überschreitet leider Ihren Kontostand",
			State:     "Failed",
			Title:     "Kauf konnte nicht durchgeführt werden",
			Operation: "-",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	/*
		TODO: VALIDATE IF THE USER HAS THE RESSOURCES TO BUY OR SELL A STOCK
	*/
	transaction := transaction.NewTransaction(userID, request.Operation, request.Operation+" "+request.StockSymbol, request.Amount, requestedStock.Price, request.Date)
	transaction.DatabaseConnection = database
	transaction.Write(true)
	createPortfolioItem(portfolio, requestedStock, currentUser, request.Amount, *totalPrice)
	/*
		TODO create portfolio_item and add to portfolio as well as reduction of balance on user
	*/
	response := TransactionResponse{
		Message:   "Kauf wird abgewickelt.. Dies kann je nach Auslastung einige Minuten dauern",
		State:     "Success",
		Title:     "Kauf abgeschlossen",
		Operation: "-",
	}
	obj, _ := json.Marshal(response)
	w.Write([]byte(obj))
}

func loadStockInstance(stockSymbol string) stock.Stock {
	stock := stock.NewStockEntry(stockSymbol, "5")
	stock.Load()
	return stock
}

//TODO implement transaction and rollback for all the queries below
func createPortfolioItem(portfolio user.Portfolio, stockInstance stock.Stock, currentUser user.User, quantity int64, totalPrice big.Float) {
	stockID := stockInstance.ID
	buyPrice := stockInstance.Price
	totalBuyPrice := totalPrice
	portfolioItem := user.PortfolioItem{
		StockID:       stockID,
		BuyPrice:      buyPrice,
		Quantity:      quantity,
		TotalBuyPrice: totalBuyPrice.String(),
	}
	portfolioItem.Write(GetDatabaseInstance())
	updatePortfolio(portfolio, totalBuyPrice, quantity, currentUser)
	connectPortfolioItemWithPortfolio(portfolio, portfolioItem, currentUser)
}

func connectPortfolioItemWithPortfolio(portfolio user.Portfolio, item user.PortfolioItem, currentUser user.User) bool {
	portfolioConnection := user.PortfolioToItem{
		PortfolioID:     portfolio.ID,
		PortfolioItemID: item.ID,
	}
	return portfolioConnection.Write(GetDatabaseInstance())
}

func AddDelayedTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Keine Parameter übergeben"))
		return
	}
	var request TransactionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
		w.Write([]byte("Invalid request format"))
		return
	}
	currentUser := user.User{
		Username: request.Username,
	}
	userID := user.GetUserIdByUsername(request.Username, GetDatabaseInstance())
	currentUser.ID = userID
	if userID <= 0 {
		logger.LogMessage(fmt.Sprintf("Eine Nutzer ID war nicht gültig | User %s", request.Username), logger.WARNING)
		return
	}
	validator := sec.NewValidator(request.AuthToken, request.Username)
	if validator.IsValidToken(jwtKey) == false {
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
	portfolio := user.LoadPortfolio(request.Username, GetDatabaseInstance())
	totalPrice := new(big.Float)
	requestedStock := loadStockInstance(request.StockSymbol)
	totalPrice, _ = totalPrice.SetString(requestedStock.Price)
	amount := float64(request.Amount)
	totalPrice = totalPrice.Mul(totalPrice, big.NewFloat(amount))
	if portfolio.Balance.Cmp(totalPrice) != 1 {
		response := TransactionResponse{
			Message:   "Dieser Kauf überschreitet leider Ihren Kontostand",
			State:     "Failed",
			Title:     "Kauf konnte nicht durchgeführt werden",
			Operation: "-",
		}
		obj, _ := json.Marshal(response)
		w.Write([]byte(obj))
		return
	}
	/*
		TODO: VALIDATE IF THE USER HAS THE RESSOURCES TO BUY OR SELL A STOCK
	*/
	transaction := transaction.NewTransaction(userID, request.Operation, request.Operation+" "+request.StockSymbol, request.Amount, request.ExpectedStockPrice, request.Date)
	transaction.DatabaseConnection = database
	transaction.Write(false)
	updatePortfolio(portfolio, *totalPrice, request.Amount, currentUser)
	response := TransactionResponse{
		Message:   "Kauf wird abgewickelt.. Dies kann je nach Auslastung einige Minuten dauern",
		State:     "Success",
		Title:     "Kauf abgeschlossen",
		Operation: "-",
	}
	obj, _ := json.Marshal(response)
	w.Write([]byte(obj))
}

func updatePortfolio(portfolio user.Portfolio, totalPrice big.Float, quantity int64, currentUser user.User) {
	newBalanceValue := portfolio.Balance.Sub(&portfolio.Balance, &totalPrice)
	newCurrentValue := portfolio.CurrentValue.Add(&portfolio.CurrentValue, &totalPrice)
	portfolio.Balance = *newBalanceValue
	portfolio.CurrentValue = *newCurrentValue
	portfolio.TotalStocks += quantity
	portfolio.Update(GetDatabaseInstance())
}
