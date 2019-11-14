package stock

type Stock struct {
	ID       int64   `json:"id"`
	Symbol   string  `json:"symbol"`
	TimeData string  `json:"timedata"`
	Price    float64 `json:"price"`
	Data     string  `json:"data"`
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
	statement, err := databaseConnection.Prepare("SELECT count(*) FROM User WHERE symbole = ? AND timedata = ?")
	defer statement.Close()
	if err != nil {
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
	statement, err := databaseConnection.Prepare("SELECT id,price,data FROM Stock WHERE symbole = ? AND timedata = ?")
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
	result.Scan(&stock.ID, &stock.Price, stock.Data)
	return true
}

func (stock *Stock) Store() {

}
