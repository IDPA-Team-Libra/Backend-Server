package database

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

//CreateTransaction creates a new transaction and writes all connections and entries to the database
func CreateTransaction(transaction transaction.Transaction, portfolio user.Portfolio, stockInstance stock.Stock, currentUser user.User, quantity int64, totalPrice big.Float, databaseConnection *sql.DB) bool {
	handler, err := databaseConnection.Begin()
	if err != nil {
		return false
	}
	if transaction.Amount <= 0 {
		return false
	}
	transaction.Symbol = stockInstance.Symbol
	newBalanceValue := portfolio.Balance.Sub(&portfolio.Balance, &totalPrice)
	if transaction.Write(true, handler, newBalanceValue.String()) == false {
		handler.Rollback()
		return false
	}
	stockID := stockInstance.ID
	buyPrice := stockInstance.Price
	totalBuyPrice := totalPrice
	portfolioItem := user.PortfolioItem{
		StockID:       stockID,
		BuyPrice:      buyPrice,
		Quantity:      quantity,
		TotalBuyPrice: totalBuyPrice.String(),
	}
	if portfolioItem.Write(handler) == false {
		handler.Rollback()
		return false
	}
	if UpdatePortfolio(portfolio, totalBuyPrice, quantity, handler) == false {
		handler.Rollback()
		return false
	}
	if connectPortfolioItemWithPortfolio(portfolio, portfolioItem, currentUser, handler) == false {
		handler.Rollback()
		return false
	}
	handler.Commit()
	return true
}

//UpdatePortfolioValues updates the values for quantity , balance and Value
func UpdatePortfolioValues(portfolio user.Portfolio, totalPrice big.Float, quantity int64) user.Portfolio {
	newBalanceValue := portfolio.Balance.Sub(&portfolio.Balance, &totalPrice)
	newCurrentValue := portfolio.CurrentValue.Add(&portfolio.CurrentValue, &totalPrice)
	portfolio.Balance = *newBalanceValue
	portfolio.CurrentValue = *newCurrentValue
	portfolio.TotalStocks += quantity
	return portfolio
}

//UpdatePortfolio updates values and writes the update to the database
func UpdatePortfolio(portfolio user.Portfolio, totalPrice big.Float, quantity int64, connection *sql.Tx) bool {
	portfolio = UpdatePortfolioValues(portfolio, totalPrice, quantity)
	return portfolio.Update(connection)
}

func connectPortfolioItemWithPortfolio(portfolio user.Portfolio, item user.PortfolioItem, currentUser user.User, databaseConnection *sql.Tx) bool {
	portfolioConnection := user.PortfolioToItem{
		PortfolioID:     portfolio.ID,
		PortfolioItemID: item.ID,
	}
	return portfolioConnection.Write(databaseConnection)
}

//UpdateUserPortfolios updates the portfolio value of all users with the currently stored stock values
func UpdateUserPortfolios(databaseConnection *sql.DB) {
	userIDs := user.GetUserIDs(databaseConnection)
	for index := range userIDs {
		res, err := databaseConnection.Query("SELECT SUM(item.quantity * (SELECT price FROM stock WHERE id = item.stock_id)) FROM portfolio_item item, portfolio_to_item port_to_item, portfolio port WHERE item.id = port_to_item.portfolio_item_id AND port_to_item.portfolio_id = port.id AND port.user_id = ?", userIDs[index])
		defer res.Close()
		if err != nil {
			return
		}
		var currentTotalPorfolioValue string
		res.Next()
		res.Scan(&currentTotalPorfolioValue)
		_, err = databaseConnection.Exec("UPDATE portfolio SET current_value = ? WHERE user_id = ?", currentTotalPorfolioValue, userIDs[index])
		if err != nil {
			fmt.Println("Error in writing update to portfolio")
		}
	}
}
