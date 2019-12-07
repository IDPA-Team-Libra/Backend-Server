package user

import (
	"fmt"

	"github.com/Liberatys/libra-back/main/stock"
)

type PortfolioItem struct {
	ID            int64
	StockID       int64
	StockName     string
	CompanyName   string
	BuyPrice      string
	Quantity      int64
	TotalBuyPrice string
	CurrentPrice  string
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

func LoadUserItems(user User) []PortfolioItem {
	statement, err := user.DatabaseConnection.Prepare("SELECT buy_price,quantity,total_buy_price,stock_id FROM portfolio_item p_i, portfolio_to_item p_i_t, Portfolio port WHERE p_i_t.portfolio_item_id = p_i.id AND p_i_t.portfolio_id = port.id AND port.user_id = ?")
	defer statement.Close()
	var items []PortfolioItem
	if err != nil {
		fmt.Println(err.Error())
		return items
	}
	fmt.Println(user.ID)
	result, err := statement.Query(user.ID)
	if err != nil {
		fmt.Println(err.Error())
		return items
	}
	defer result.Close()
	for result.Next() {
		var item PortfolioItem
		result.Scan(&item.BuyPrice, &item.Quantity, &item.TotalBuyPrice, &item.StockID)
		stockInstance := stock.Stock{}
		symbol := stockInstance.GetSymbolByID(item.StockID)
		stockInstance = stock.NewStockEntry(symbol, "5")
		stockInstance.Load()
		item.CurrentPrice = stockInstance.Price
		item.CompanyName = stockInstance.Company
		item.StockName = stockInstance.Symbol
		items = append(items, item)
	}
	return items
}
