package database

import (
	"database/sql"
	"fmt"
	"math/big"
	"strings"

	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/mail"
	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

//Outcome the result of a transaction execution
type Outcome struct {
	Operation string
	Sucess    bool
	Symbol    string
	Message   string
	Amount    int64
}

//StartBatchProcess collect all delayed transaction for the current date and execute them
func StartBatchProcess(databaseConnection *sql.DB) {
	transactions := LoadDelayedTransactions(databaseConnection)
	var operations map[int64][]Outcome
	operations = make(map[int64][]Outcome)
	for _, value := range transactions {
		_, contains := operations[value.UserID]
		outcome := true
		message := ""
		if value.Action == "buy" {
			outcome, message = BatchBuyTransactions(value, databaseConnection)
		} else {
			outcome, message = BatchSellTransactions(value, databaseConnection)
		}
		if contains == false {
			operations[value.UserID] = []Outcome{
				Outcome{
					Sucess:    outcome,
					Operation: value.Action,
					Symbol:    value.Symbol,
					Message:   message,
					Amount:    value.Amount,
				},
			}
		} else {
			operations[value.UserID] = append(operations[value.UserID], Outcome{
				Sucess:    outcome,
				Operation: value.Action,
				Symbol:    value.Symbol,
				Message:   message,
				Amount:    value.Amount,
			})
		}
	}
	SendUpdatesPerUser(operations, databaseConnection)
}

//SendUpdatesPerUser sends an email with the outcome of his transaction execution
func SendUpdatesPerUser(mapping map[int64][]Outcome, databaseConnection *sql.DB) {
	for key, value := range mapping {
		userInstance := user.User{
			ID: key,
		}
		email := userInstance.GetUserMail(databaseConnection)
		totalOperations := len(value)
		var result []string
		for index := range value {
			val := value[index]
			sucessState := ""
			if val.Sucess == true {
				sucessState = "Durchgeführt"
			} else {
				sucessState = "Fehlgeschlagen"
			}
			currentStat := fmt.Sprintf("%s : %s => Anzahl: %d ==> %s", val.Operation, val.Symbol, val.Amount, sucessState)
			if val.Sucess == false {
				currentStat = fmt.Sprintf("%s | %s => Anzahl: %d", currentStat, val.Message, val.Amount)
			}
			result = append(result, currentStat)
		}
		mailer := mail.NewMail(email)
		mailer.ApplyConfiguration(mail.LoadMailConfiguration())
		go mailer.SendDelayedTransactionEmail(totalOperations, result)
	}
}

//LoadDelayedTransactions loads delayed transactions
func LoadDelayedTransactions(databaseConnection *sql.DB) []transaction.Transaction {
	transaction := transaction.Transaction{}
	return transaction.LoadTransactionsByProcessState(-1, databaseConnection, false)
}

//BatchBuyTransactions handles delayed sell transactions
func BatchBuyTransactions(transaction transaction.Transaction, conn *sql.DB) (bool, string) {
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
	addition := new(big.Float)
	addition, _ = addition.SetString("50")
	totalPrice = totalPrice.Add(totalPrice, addition)
	// can only execute enough money is around
	transaction.Remove(handler)
	handler.Commit()
	if portfolio.Balance.Cmp(totalPrice) == 2 {
		logger.LogMessage(fmt.Sprintf("Was not able to complete delayed buy for User: %s due to insufficient funds", currentUser.Username), logger.WARNING)
		return false, "Zu wenig Geld auf der Bank"
	}
	if CreateTransaction(transaction, portfolio, requestedStock, currentUser, transaction.Amount, *totalPrice, conn) == false {
		return false, "Server-Problem"
	}
	return true, "Alles Super"
}

//ExtractStockNameWithTrim extracts the stockname from a field
func ExtractStockNameWithTrim(description string) string {
	return strings.TrimSpace(ExtractStockName(description))
}

//ExtractStockName extract the stockname from a field
func ExtractStockName(description string) string {
	parts := strings.Split(description, " ")
	if len(parts) == 1 {
		return parts[0]
	}
	return parts[1]
}

//BatchSellTransactions handles the delayed buy transactions
func BatchSellTransactions(transaction transaction.Transaction, conn *sql.DB) (bool, string) {
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
		return false, "Zu wenig Aktien"
	}
	handler, err := conn.Begin()
	if err != nil {
		handler.Rollback()
		fmt.Println(err.Error())
		return false, "Server-Problem"
	}
	changedItems := SubtractStocksFromTotalAmount(items, requestCount)
	UpdateOrDeleteStocks(changedItems, handler)
	userID := user.GetUserIDByUsername(currentUser.Username, conn)
	portfolio := user.LoadPortfolio(userID, conn)
	portfolio.TotalStocks -= transaction.Amount
	s := fmt.Sprintf("%f", float64(transaction.Amount))
	additionalBalance := MultiplyString(s, requestedStock.Price)
	portfolio.Balance = *portfolio.Balance.Add(&portfolio.Balance, additionalBalance)
	subtraction := new(big.Float)
	subtraction, _ = subtraction.SetString("50")
	additionalBalance = additionalBalance.Sub(additionalBalance, subtraction)
	portfolio.CurrentValue = *portfolio.CurrentValue.Sub(&portfolio.CurrentValue, additionalBalance)
	if transaction.Write(true, handler, portfolio.Balance.String()) == false {
		handler.Rollback()
		return false, "Server-Problem"
	}
	if portfolio.Update(handler) == false {
		handler.Rollback()
		return false, "Server-Problem"
	}
	logger.LogMessage(fmt.Sprintf("Executed sell for User: %s", currentUser.Username), logger.WARNING)
	handler.Commit()
	return true, "Alles Super"
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
