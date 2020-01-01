package stock

import (
	"fmt"
)

//Stock struct holding informaiton about stocks in the database
type Stock struct {
	ID       int64  `json:"id"`
	Company  string `json:"company"`
	Symbol   string `json:"symbol"`
	Price    string `json:"price"`
	TimeData string `json:"timedata"`
	Data     string `json:"data"`
}

//NewStockEntry creates a new entry with symbol and timedata
func NewStockEntry(Symbol string, TimeData string) Stock {
	stock := Stock{
		Symbol:   Symbol,
		TimeData: TimeData,
	}
	return stock
}

//Load retreaves stock information from the database
func (stock *Stock) Load() bool {
	databaseConnection := database
	statement, err := databaseConnection.Prepare("SELECT id,data,price,company FROM stock WHERE symbol = ? AND timedata = ?")
	if err != nil {
		return false
	}
	defer statement.Close()
	result, err := statement.Query(stock.Symbol, stock.TimeData)
	if err != nil {
		return false
	}
	defer result.Close()
	result.Next()
	result.Scan(&stock.ID, &stock.Data, &stock.Price, &stock.Company)
	return true
}

//Update updates the stock values in the database
func (stock *Stock) Update() bool {
	databaseConnection := database
	statement, err := databaseConnection.Prepare("UPDATE stock SET data = ?,price = ?,company = ? WHERE id = ? AND timeData = ?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = statement.Exec(stock.Data, stock.Price, stock.Company, stock.ID, stock.TimeData)
	if err != nil {
		return false
	}
	return true
}

//GetSymbolByID returns a stock symbol by comparing id and timedata
func (stock *Stock) GetSymbolByID(id int64) string {
	databaseConnection := database
	statement, err := databaseConnection.Prepare("SELECT symbol FROM stock WHERE id = ? AND timedata = ?")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer statement.Close()
	result, err := statement.Query(id, "5")
	if err != nil {
		return ""
	}
	defer result.Close()
	result.Next()
	result.Scan(&stock.Symbol)
	return stock.Symbol
}

//LoadStockInstance loads a stock by a given symbol
func LoadStockInstance(stockSymbol string) Stock {
	stock := NewStockEntry(stockSymbol, "5")
	stock.Load()
	return stock
}
