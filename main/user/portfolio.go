package user

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/Liberatys/libra-back/main/logger"
)

type Portfolio struct {
	ID           int64           `json:"id"`
	Balance      big.Float       `json:"balance"`
	CurrentValue big.Float       `json:"currentValue"`
	TotalStocks  int64           `json:"totalStocks"`
	StartCapital big.Float       `json:"startCapital"`
	Items        []PortfolioItem `json:"items"`
}
type StubReader struct {
	CurrentValue string
	Balance      string
	TotalStocks  int64
	StartCapital string
	ID           int64
}

func LoadPortfolio(username string, connection *sql.DB) Portfolio {
	statement, err := connection.Prepare("SELECT id,current_value, total_stocks, start_capital,balance FROM Portfolio WHERE user_id = ?")
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
		statement.Close()
	}
	defer statement.Close()
	result, err := statement.Query(GetUserIdByUsername(username, connection))
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
		result.Close()
		return Portfolio{}
	}
	defer result.Close()
	var reader StubReader
	for result.Next() {
		err = result.Scan(&reader.ID, &reader.CurrentValue, &reader.TotalStocks, &reader.StartCapital, &reader.Balance)
		if err != nil {
			logger.LogMessage(err.Error(), logger.WARNING)
		}
	}
	return ConvertStub(reader)
}

func ConvertStub(reader StubReader) Portfolio {
	portfolio := Portfolio{}
	portfolio.CurrentValue = convertStringToFloat(reader.CurrentValue)
	portfolio.StartCapital = convertStringToFloat(reader.StartCapital)
	portfolio.Balance = convertStringToFloat(reader.Balance)
	portfolio.ID = reader.ID
	portfolio.TotalStocks = reader.TotalStocks
	return portfolio
}

func convertStringToFloat(value string) big.Float {
	currentVal := new(big.Float)
	currentVal.SetString(value)
	return *currentVal
}

func (portfolio *Portfolio) Write(userID int64, connection *sql.DB, startCapital float64) bool {
	portfolio.ID = userID
	portfolio.CurrentValue = *big.NewFloat(0.0)
	portfolio.TotalStocks = 0
	portfolio.Balance = *big.NewFloat(startCapital)
	portfolio.StartCapital = *big.NewFloat(startCapital)
	statement, err := connection.Prepare("INSERT INTO Portfolio(user_id, current_value,total_stocks,start_capital,balance) VALUES(?,?,0,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(userID, 0.0, startCapital, startCapital)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (portfolio *Portfolio) Update(connection *sql.Tx) bool {
	statement, err := connection.Prepare("UPDATE Portfolio SET balance = ?, current_value = ?, total_stocks = ? WHERE id = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(portfolio.Balance.String(), portfolio.CurrentValue.String(), portfolio.TotalStocks, portfolio.ID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (portfolio *Portfolio) AddItem(portfolioItem PortfolioItem, connection *sql.DB) bool {
	statement, err := connection.Prepare("INSERT INTO portfolio_to_item(portfolio_id,portfolio_item_id) VALUES(?,?)")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(portfolio.ID, portfolioItem.ID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
