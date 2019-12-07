package user

import (
	"fmt"
	"math/big"
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

func LoadPortfolio(user User) Portfolio {
	statement, err := user.DatabaseConnection.Prepare("SELECT id,current_value, total_stocks, start_capital,balance FROM Portfolio WHERE user_id = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := statement.Query(user.GetUserIdByUsername(user.Username))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	result.Next()
	var reader StubReader
	err = result.Scan(&reader.ID, &reader.CurrentValue, &reader.TotalStocks, &reader.StartCapital, &reader.Balance)
	fmt.Println(reader)
	if err != nil {
		fmt.Println(err.Error())
	}
	return ConvertStub(reader)
}

func ConvertStub(reader StubReader) Portfolio {
	portfolio := Portfolio{}
	portfolio.CurrentValue = convertStringToFloat(reader.CurrentValue)
	portfolio.StartCapital = convertStringToFloat(reader.StartCapital)
	portfolio.Balance = convertStringToFloat(reader.Balance)
	portfolio.ID = reader.ID
	return portfolio
}

func convertStringToFloat(value string) big.Float {
	currentVal := new(big.Float)
	currentVal.SetString(value)
	return *currentVal
}

func (portfolio *Portfolio) Write(userID int64, user User, startCapital float64) bool {
	portfolio.ID = userID
	portfolio.CurrentValue = *big.NewFloat(startCapital)
	portfolio.TotalStocks = 0
	portfolio.StartCapital = *big.NewFloat(startCapital)
	statement, err := user.DatabaseConnection.Prepare("INSERT INTO Portfolio(user_id, current_value,total_stocks,start_capital,balance) VALUES(?,?,0,?,?)")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(userID, startCapital, startCapital, startCapital)
	if err != nil {
		return false
	}
	return true
}

func (portfolio *Portfolio) Update(user User) bool {
	statement, err := user.DatabaseConnection.Prepare("UPDATE Portfolio SET balance = ?, current_value = ?, total_stocks = ? WHERE id = ?")
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

func (portfolio *Portfolio) AddItem(portfolioItem PortfolioItem, user User) bool {
	statement, err := user.DatabaseConnection.Prepare("INSERT INTO portfolio_to_item(portfolio_id,portfolio_item_id) VALUES(?,?)")
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

func IncrementPortfolioTotalStockAndPrice() {

}

func (portfolio *Portfolio) RemoveItem() {

}

func (portfolio *Portfolio) StockAlreadyPresent() {

}
