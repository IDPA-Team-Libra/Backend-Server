package user

import (
	"database/sql"
	"math/big"

	"github.com/Liberatys/libra-back/main/logger"
)

//Portfolio a structure to hold information representing a portfolio in the database
type Portfolio struct {
	ID           int64           `json:"id"`
	Balance      big.Float       `json:"balance"`
	CurrentValue big.Float       `json:"currentValue"`
	TotalStocks  int64           `json:"totalStocks"`
	StartCapital big.Float       `json:"startCapital"`
	Items        []PortfolioItem `json:"items"`
}

//StubReader temporary struct for parsing a portfolio from the database
type StubReader struct {
	CurrentValue string
	Balance      string
	TotalStocks  int64
	StartCapital string
	ID           int64
}

//LoadPortfolio loads portfolio-entries for a given user
func LoadPortfolio(userID int64, connection *sql.DB) Portfolio {
	statement, err := connection.Prepare("SELECT id,current_value, total_stocks, start_capital,balance FROM portfolio WHERE user_id = ?")
	defer statement.Close()
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
	}
	result, err := statement.Query(userID)
	defer result.Close()
	if err != nil {
		logger.LogMessage(err.Error(), logger.WARNING)
		return Portfolio{}
	}
	var reader StubReader
	for result.Next() {
		err = result.Scan(&reader.ID, &reader.CurrentValue, &reader.TotalStocks, &reader.StartCapital, &reader.Balance)
		if err != nil {
			logger.LogMessage(err.Error(), logger.WARNING)
		}
	}
	return ConvertStub(reader)
}

//ConvertStub reads values and converts the strings to big.Floats
func ConvertStub(reader StubReader) Portfolio {
	portfolio := Portfolio{}
	portfolio.CurrentValue = ConvertStringToBigFloat(reader.CurrentValue)
	portfolio.StartCapital = ConvertStringToBigFloat(reader.StartCapital)
	portfolio.Balance = ConvertStringToBigFloat(reader.Balance)
	portfolio.ID = reader.ID
	portfolio.TotalStocks = reader.TotalStocks
	return portfolio
}

//ConvertStringToBigFloat convert a string to a big.Float
func ConvertStringToBigFloat(value string) big.Float {
	currentVal := new(big.Float)
	currentVal.SetString(value)
	return *currentVal
}

func (portfolio *Portfolio) Write(userID int64, connection *sql.DB, startCapital float64) bool {
	portfolio.ID = userID
	portfolio.SetDefaultValuesForPortfolio(startCapital)
	statement, err := connection.Prepare("INSERT INTO portfolio(user_id, current_value,total_stocks,start_capital,balance) VALUES(?,?,0,?,?)")
	if err != nil {
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(userID, 0.0, startCapital, startCapital)
	if err != nil {
		return false
	}
	return true
}

//SetDefaultValuesForPortfolio sets the default values for a portfolio
func (portfolio *Portfolio) SetDefaultValuesForPortfolio(startCapital float64) {
	portfolio.CurrentValue = *big.NewFloat(0.0)
	portfolio.TotalStocks = 0
	portfolio.Balance = *big.NewFloat(startCapital)
	portfolio.StartCapital = *big.NewFloat(startCapital)
}

//Update updates an existing portfolio value in the database
func (portfolio *Portfolio) Update(connection *sql.Tx) bool {
	statement, err := connection.Prepare("UPDATE portfolio SET balance = ?, current_value = ?, total_stocks = ? WHERE id = ?")
	defer statement.Close()
	if err != nil {
		return false
	}
	_, err = statement.Exec(portfolio.Balance.String(), portfolio.CurrentValue.String(), portfolio.TotalStocks, portfolio.ID)
	if err != nil {
		return false
	}
	return true
}

//AddItem adds a connection between a portfolio and a portfolioItem
func (portfolio *Portfolio) AddItem(portfolioItem PortfolioItem, connection *sql.DB) bool {
	statement, err := connection.Prepare("INSERT INTO portfolio_to_item(portfolio_id,portfolio_item_id) VALUES(?,?)")
	defer statement.Close()
	if err != nil {
		return false
	}
	_, err = statement.Exec(portfolio.ID, portfolioItem.ID)
	if err != nil {
		return false
	}
	return true
}
