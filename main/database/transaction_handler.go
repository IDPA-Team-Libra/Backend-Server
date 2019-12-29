package database

import (
	"database/sql"
	"math/big"

	"github.com/Liberatys/libra-back/main/stock"
	"github.com/Liberatys/libra-back/main/transaction"
	"github.com/Liberatys/libra-back/main/user"
)

func CreateTransaction(transaction transaction.Transaction, portfolio user.Portfolio, stockInstance stock.Stock, currentUser user.User, quantity int64, totalPrice big.Float, databaseConnection *sql.DB) bool {
	handler, err := databaseConnection.Begin()
	if err != nil {
		return false
	}
	if transaction.Write(true, handler) == false {
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
	if UpdatePortfolio(portfolio, totalBuyPrice, quantity, currentUser, handler) == false {
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

func UpdatePortfolio(portfolio user.Portfolio, totalPrice big.Float, quantity int64, currentUser user.User, connection *sql.Tx) bool {
	newBalanceValue := portfolio.Balance.Sub(&portfolio.Balance, &totalPrice)
	newCurrentValue := portfolio.CurrentValue.Add(&portfolio.CurrentValue, &totalPrice)
	portfolio.Balance = *newBalanceValue
	portfolio.CurrentValue = *newCurrentValue
	portfolio.TotalStocks += quantity
	return portfolio.Update(connection)
}

func connectPortfolioItemWithPortfolio(portfolio user.Portfolio, item user.PortfolioItem, currentUser user.User, db_connection *sql.Tx) bool {
	portfolioConnection := user.PortfolioToItem{
		PortfolioID:     portfolio.ID,
		PortfolioItemID: item.ID,
	}
	return portfolioConnection.Write(db_connection)
}
