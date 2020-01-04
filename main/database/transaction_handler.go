package database

import (
	"database/sql"
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
