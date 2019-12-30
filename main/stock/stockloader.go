package stock

import (
	"fmt"

	av "github.com/cmckee-dev/go-alpha-vantage"
)

func LoadStocksForRoute(timeSeries string) []Stock {
	database_connection := database
	var stocks []Stock
	statement, err := database_connection.Prepare("SELECT id,symbol,timeData FROM stock where timeData = ? AND price > 0")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := statement.Query(timeSeries)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var stock Stock
		result.Scan(&stock.ID, &stock.Symbol, &stock.TimeData)
		stocks = append(stocks, stock)
	}
	return stocks
}

func LoadAllStockSymbols(timeSeries string) []Stock {
	database_connection := database
	var stocks []Stock
	statement, err := database_connection.Prepare("SELECT id,symbol,timeData FROM stock where timeData = ?")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err := statement.Query(timeSeries)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var stock Stock
		result.Scan(&stock.ID, &stock.Symbol, &stock.TimeData)
		stocks = append(stocks, stock)
	}
	return stocks
}

func ConvertTimeSeries(series string) av.TimeInterval {
	switch series {
	case "5":
		return av.TimeIntervalFiveMinute
	case "15":
		return av.TimeIntervalFifteenMinute
	case "30":
		return av.TimeIntervalThirtyMinute
	case "60":
		return av.TimeIntervalSixtyMinute
	default:
		return av.TimeIntervalFiveMinute
	}
}
