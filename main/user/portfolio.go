package user

import (
	"fmt"
	"math/big"
)

type Portfolio struct {
	ID           int64     `json:"id"`
	CurrentValue big.Float `json:"currentValue"`
	TotalStocks  int64     `json:"totalStocks"`
	StartCapital big.Float `json:"startCapital"`
}

func LoadPortfolio(user User) Portfolio {
	statement, err := user.DatabaseConnection.Prepare("SELECT current_value, total_stocks, start_capital FROM Portfolio WHERE user_id = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := statement.Query(user.GetUserIdByUsername(user.Username))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	var portfolio Portfolio
	result.Next()
	result.Scan(&portfolio.CurrentValue, &portfolio.TotalStocks, &portfolio.StartCapital)
	return portfolio
}

func (portfolio *Portfolio) Create(userID int64, user User, startCapital float64) bool {
	portfolio.ID = userID
	portfolio.CurrentValue = *big.NewFloat(startCapital)
	portfolio.TotalStocks = 0
	portfolio.StartCapital = *big.NewFloat(startCapital)
	statement, err := user.DatabaseConnection.Prepare("INSERT INTO Portfolio(user_id, current_value,total_stocks,start_capital) VALUES(?,?,0,?)")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(userID, startCapital, startCapital)
	if err != nil {
		return false
	}
	return true
}

func (portfolio *Portfolio) Write() {

}

func (portfolio *Portfolio) AddItem() {

}

func (portfolio *Portfolio) RemoveItem() {

}

func (portfolio *Portfolio) StockAlreadyPresent() {

}
