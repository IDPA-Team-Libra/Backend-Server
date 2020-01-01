package user

import (
	"database/sql"
	"fmt"

	"github.com/Liberatys/libra-back/main/stock"
)

//PortfolioItem functions as a holder for values relevant to the porfolio and stock
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

//Write writes the portfolio item to the database
func (item *PortfolioItem) Write(connection *sql.Tx) bool {
	statement, err := connection.Prepare("INSERT INTO portfolio_item(stock_id,buy_price, quantity, total_buy_price,buy_date_time) VALUES(?,?,?,?,NOW())")
	defer statement.Close()
	if err != nil {
		return false
	}
	resp, err := statement.Exec(item.StockID, item.BuyPrice, item.Quantity, item.TotalBuyPrice)
	if err != nil {
		return false
	}
	value, err := resp.LastInsertId()
	if err != nil {
		return false
	}
	item.ID = value
	return true
}

func LoadUserItems(userID int64, Symbol string, connection *sql.DB) []PortfolioItem {
	statement, err := connection.Prepare("SELECT p_i.id, buy_price,quantity,total_buy_price,stock_id FROM portfolio_item p_i, portfolio_to_item p_i_t, Portfolio port WHERE p_i_t.portfolio_item_id = p_i.id AND p_i_t.portfolio_id = port.id AND port.user_id = ?")
	defer statement.Close()
	var items []PortfolioItem
	if err != nil {
		return items
	}
	result, err := statement.Query(userID)
	if err != nil {
		return items
	}
	defer result.Close()
	for result.Next() {
		var item PortfolioItem
		result.Scan(&item.ID, &item.BuyPrice, &item.Quantity, &item.TotalBuyPrice, &item.StockID)
		stockInstance := stock.Stock{}
		symbol := stockInstance.GetSymbolByID(item.StockID)
		stockInstance = stock.NewStockEntry(symbol, "5")
		stockInstance.Load()
		if Symbol != "*" {
			if Symbol != stockInstance.Symbol {
				continue
			}
		}
		item.CurrentPrice = stockInstance.Price
		item.CompanyName = stockInstance.Company
		item.StockName = stockInstance.Symbol
		items = append(items, item)
	}
	return items
}

func (item *PortfolioItem) Update(connetion *sql.Tx) bool {
	statement, err := connetion.Prepare("UPDATE portfolio_item SET quantity = ? WHERE id = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(item.Quantity, item.ID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (item *PortfolioItem) RemoveItemAndConnection(connection *sql.Tx) bool {
	if item.RemoveConnection(connection) {
		return item.Remove(connection)
	}
	return false
}

func (item *PortfolioItem) RemoveConnection(databaseConnection *sql.Tx) bool {
	return RemoveEntryByQuery(databaseConnection, "DELETE FROM portfolio_to_item WHERE portfolio_item_id = ?", item.ID)
}

func (item *PortfolioItem) Remove(databaseConnection *sql.Tx) bool {
	return RemoveEntryByQuery(databaseConnection, "DELETE FROM portfolio_item WHERE id = ?", item.ID)
}

func RemoveEntryByQuery(databaseConnection *sql.Tx, query string, parameters ...interface{}) bool {
	statement, err := databaseConnection.Prepare(query)
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(parameters...)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
