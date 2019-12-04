package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

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
}

type FutureTransactionOption struct {
	TransactionRequest TransactionRequest `json:"transactionRequest"`
	SetDate            string             `json:"setDate"`
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
	currentUser.SetDatabaseConnection(database)
	userID := currentUser.GetUserIdByUsername(request.Username)
	currentUser.ID = userID
	if userID <= 0 {
		fmt.Println("Invalud userID")
		return
	}
	portfolio := user.LoadPortfolio(currentUser)
	fmt.Println(portfolio)
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
	transaction := transaction.NewTransaction(userID, request.Operation, "", request.Amount, request.ExpectedStockPrice, request.Date)
	transaction.DatabaseConnection = database
	transaction.Write()
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
	portfolioItem.Write(currentUser)
	newBalanceValue := portfolio.Balance.Sub(&portfolio.Balance, &totalBuyPrice)
	portfolio.Balance = *newBalanceValue
	portfolio.TotalStocks += quantity
	portfolio.Update(currentUser)
}

func AddDelayedTransaction(w http.ResponseWriter, r *http.Request) {
}
