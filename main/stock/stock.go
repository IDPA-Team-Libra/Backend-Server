package stock

import (
	"fmt"
)

type Stock struct {
	ID       int64  `json:"id"`
	Company  string `json:"company"`
	Symbol   string `json:"symbol"`
	Price    string `json:"price"`
	TimeData string `json:"timedata"`
	Data     string `json:"data"`
}

func NewStockEntry(Symbol string, TimeData string) Stock {
	stock := Stock{
		Symbol:   Symbol,
		TimeData: TimeData,
	}
	return stock
}

func (stock *Stock) IsPresent() bool {
	databaseConnection := database
	statement, err := databaseConnection.Prepare("SELECT count(*) FROM stock WHERE symbol = ? AND timeData = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	result, err := statement.Query(stock.Symbol, stock.TimeData)
	if err != nil {
	}
	defer result.Close()
	var returnedCounter int
	result.Next()
	result.Scan(&returnedCounter)
	if returnedCounter == 0 {
		return false
	}
	return true
}

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

func (stock *Stock) Store() bool {
	databaseConnection := database
	if stock.IsPresent() {
		statement, err := databaseConnection.Prepare("UPDATE stock SET data = ?,price = ?,company = ? WHERE id = ? AND timeData = ?")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		_, err = statement.Exec(stock.Data, stock.Price, stock.Company, stock.ID, stock.TimeData)
		if err != nil {
			return false
		}
	} else {
		statement, err := databaseConnection.Prepare("INSERT INTO stock(symbol,company,timeData,data,last_query,price) VALUES(?,?,?,?,Now(),?)")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		_, err = statement.Exec(stock.Symbol, stock.Company, stock.TimeData, stock.Data, stock.Price)
		if err != nil {
			return false
		}
	}
	return true
}

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

func LoadStockInstance(stockSymbol string) Stock {
	stock := NewStockEntry(stockSymbol, "5")
	stock.Load()
	return stock
}
