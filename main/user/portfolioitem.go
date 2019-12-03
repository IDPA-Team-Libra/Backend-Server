package user

import "fmt"

type PortfolioItem struct {
	ID            int64
	StockID       int64
	BuyPrice      string
	Quantity      int64
	TotalBuyPrice string
}

func (item *PortfolioItem) Write(user User) bool {
	statement, err := user.DatabaseConnection.Prepare("INSERT INTO portfolio_item(stock_id,buy_price, quantity, total_buy_price,buy_date_time) VALUES(?,?,?,?,NOW())")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	resp, err := statement.Exec(item.StockID, item.BuyPrice, item.Quantity, item.TotalBuyPrice)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	value, err := resp.LastInsertId()
	if err != nil {
		return false
	}
	item.ID = value
	return true
}

func (item *PortfolioItem) Load() {

}
