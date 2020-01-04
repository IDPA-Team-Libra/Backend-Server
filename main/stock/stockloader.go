package stock

import (
	"fmt"

	av "github.com/Liberatys/go-alpha-vantage"
)

//LoadStocksForRoute loads all stocks for the stockRoute
func LoadStocksForRoute(timeSeries string) []Stock {
	databaseConnection := database
	var stocks []Stock
	statement, err := databaseConnection.Prepare("SELECT id,symbol,timeData FROM stock where timeData = ? AND price > 0")
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

//LoadAllStockSymbols gets all stocks from the database that have a given timeSeries [1,5,15,30,60]
func LoadAllStockSymbols(timeSeries string) []Stock {
	databaseConnection := database
	var stocks []Stock
	statement, err := databaseConnection.Prepare("SELECT id,symbol,timeData FROM stock where timeData = ?")
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

//ConvertTimeSeries converts a string into an av datatype
func ConvertTimeSeries(series string) (av.TimeInterval, bool) {
	switch series {
	case "1":
		return av.TimeIntervalOneMinute, true
	case "5":
		return av.TimeIntervalFiveMinute, true
	case "15":
		return av.TimeIntervalFifteenMinute, true
	case "30":
		return av.TimeIntervalThirtyMinute, true
	case "60":
		return av.TimeIntervalSixtyMinute, true
	default:
		return av.TimeIntervalFiveMinute, false
	}
}

//ConvertTimeCap a method to convert a string to an av.TimeSeries
func ConvertTimeCap(cap string) av.TimeSeries {
	switch cap {
	case "Daily":
		return av.TimeSeriesDaily
	case "Weekly":
		return av.TimeSeriesWeekly
	case "Monthly":
		return av.TimeSeriesMonthly
	default:
		return av.TimeSeriesDaily
	}
}
