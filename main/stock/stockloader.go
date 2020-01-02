package stock

import (
	"fmt"

	av "github.com/Liberatys/go-alpha-vantage"
)

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
