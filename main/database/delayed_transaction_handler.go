package database

import (
	"database/sql"
	"fmt"
	"math/big"
	"strings"

	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

func StartBatchProcess(databaseConnection *sql.DB) {
	transactions := LoadDelayedTransactions(databaseConnection)
	for _, value := range transactions {
		if value.Action == "buy" {
			BatchBuyTransactions(value, databaseConnection)
		} else {
			BatchSellTransactions(value, databaseConnection)
		}
	}
}

func LoadDelayedTransactions(databaseConnection *sql.DB) []transaction.Transaction {
	transaction := transaction.Transaction{}
	return transaction.LoadTransactionsByProcessState(-1, databaseConnection, false)
}

func BatchBuyTransactions(transaction transaction.Transaction, conn *sql.DB) bool {
	currentUser := user.User{
		Username: user.GetUsernameByID(transaction.UserID, conn),
	}
	userID := user.GetUserIDByUsername(currentUser.Username, conn)
	portfolio := user.LoadPortfolio(userID, conn)
	totalPrice := new(big.Float)
	requestedStock := stock.LoadStockInstance(ExtractStockNameWithTrim(transaction.Symbol))
	totalPrice, _ = totalPrice.SetString(requestedStock.Price)
	amount := float64(transaction.Amount)
	totalPrice = totalPrice.Mul(totalPrice, big.NewFloat(amount))
	expectedTransactionValue := new(big.Float)
	expectedTransactionValue, _ = expectedTransactionValue.SetString(transaction.Value)
	handler, _ := conn.Begin()
	// can only execute enough money is around
	transaction.Remove(handler)
	handler.Commit()
	if portfolio.Balance.Cmp(totalPrice) == 2 {
		logger.LogMessage(fmt.Sprintf("Was not able to complete delayed buy for User: %s due to insufficient funds", currentUser.Username), logger.WARNING)
		return false
	}
	if CreateTransaction(transaction, portfolio, requestedStock, currentUser, transaction.Amount, *totalPrice, conn) == false {
		return false
	}
	return true
}

func ExtractStockNameWithTrim(description string) string {
	return strings.TrimSpace(ExtractStockName(description))
}

func ExtractStockName(description string) string {
	parts := strings.Split(description, " ")
	if len(parts) == 1 {
		return parts[0]
	}
	return parts[1]
}

func BatchSellTransactions(transaction transaction.Transaction, conn *sql.DB) bool {
	currentUser := user.User{
		Username: user.GetUsernameByID(transaction.UserID, conn),
	}
	stockSymbol := ExtractStockNameWithTrim(transaction.Symbol)
	userInstance := user.CreateUserInstance(currentUser.Username, "", "")
	userInstance.ID = user.GetUserIDByUsername(currentUser.Username, conn)
	requestedStock := stock.LoadStockInstance(stockSymbol)
	items := user.LoadUserItems(userInstance.ID, stockSymbol, conn)
	totalStockQuantity := CalculateTotalStocks(items)
	requestCount := transaction.Amount
	handler, _ := conn.Begin()
	transaction.Remove(handler)
	handler.Commit()
	if totalStockQuantity < requestCount {
		logger.LogMessage(fmt.Sprintf("Was not able to execute delayed sell for User: %s", currentUser.Username), logger.WARNING)
		return false
	}
	handler, err := conn.Begin()
	if err != nil {
		handler.Rollback()
		fmt.Println(err.Error())
		return false
	}
	changedItems := SubtractStocksFromTotalAmount(items, requestCount)
	UpdateOrDeleteStocks(changedItems, handler)
	if transaction.Write(true, handler) == false {
		handler.Rollback()
		return false
	}
	userID := user.GetUserIDByUsername(currentUser.Username, conn)
	portfolio := user.LoadPortfolio(userID, conn)
	portfolio.TotalStocks -= transaction.Amount
	s := fmt.Sprintf("%f", float64(transaction.Amount))
	additionalBalance := MultiplyString(s, requestedStock.Price)
	portfolio.Balance = *portfolio.Balance.Add(&portfolio.Balance, additionalBalance)
	portfolio.CurrentValue = *portfolio.CurrentValue.Sub(&portfolio.CurrentValue, additionalBalance)
	if portfolio.Update(handler) == false {
		handler.Rollback()
		return false
	}
	logger.LogMessage(fmt.Sprintf("Executed sell for User: %s", currentUser.Username), logger.WARNING)
	handler.Commit()
	return true
}

//MultiplyString multiplies two strings and returns it as *big.Float
func MultiplyString(first, second string) *big.Float {
	firstFloat := new(big.Float)
	firstFloat.SetString(first)
	secondFloat := new(big.Float)
	firstFloat, _ = firstFloat.SetString(first)
	secondFloat, _ = secondFloat.SetString(second)
	return firstFloat.Mul(firstFloat, secondFloat)
}

//SubtractStocksFromTotalAmount subtracts a count from a slice of PortfolioItems
func SubtractStocksFromTotalAmount(items []user.PortfolioItem, requestCount int64) []user.PortfolioItem {
	var changedItems []user.PortfolioItem
	for index := range items {
		if requestCount > 0 {
			quantity := items[index].Quantity
			if quantity <= requestCount {
				requestCount -= quantity
				items[index].Quantity = 0
				changedItems = append(changedItems, items[index])
			} else {
				items[index].Quantity -= requestCount
				requestCount = 0
				changedItems = append(changedItems, items[index])
				break
			}
		} else {
			break
		}
	}
	return changedItems
}

//UpdateOrDeleteStocks updates or deletes stocks based on the Quntity value
func UpdateOrDeleteStocks(items []user.PortfolioItem, conn *sql.Tx) bool {
	for index := range items {
		quantity := items[index].Quantity
		if quantity > 0 {
			if items[index].Update(conn) == false {
				return false
			}
		} else {
			if items[index].RemoveItemAndConnection(conn) == false {
				return false
			}
		}
	}
	return true
}

//CalculateTotalStocks calculate the total Quantity for all stocks of a given symbol
func CalculateTotalStocks(items []user.PortfolioItem) int64 {
	var counter int64
	for index := range items {
		counter += items[index].Quantity
	}
	return counter
}
