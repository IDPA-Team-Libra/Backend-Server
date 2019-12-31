package apiconnection

import (
	"fmt"
	"sync"
	"time"

	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/stock"
	av "github.com/cmckee-dev/go-alpha-vantage"
)

//TODO hide api_key
const APIKEY = "CG96DXD2YPARDLMX"

func GetStockDataForSymbol(recovered_stock stock.Stock, interval av.TimeInterval) (stock.Stock, bool) {
	client := av.NewClient(APIKEY)
	result, err := client.StockTimeSeriesIntraday(interval, recovered_stock.Symbol)
	if err != nil {
	}
	if err != nil {
		return stock.Stock{}, false
	}
	if len(result) == 0 {
		return recovered_stock, false
	}
	price := fmt.Sprintf("%.3f", result[len(result)-1].Close)
	recovered_stock.Price = price
	return recovered_stock, true
}

//const APIKEY = "F1HqA-xHQe7tzNYtFf26"
//func GetStockDataForSymbol(recovered_stock stock.Stock, interval av.TimeInterval) (stock.Stock, bool) {
//	quandl.APIKey = APIKEY
//	data, err := quandl.GetSymbol("WIKI/"+recovered_stock.Symbol, nil)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	if len(data.Data) == 0 {
//		return recovered_stock, false
//	}
//	value := fmt.Sprintf("%v", data.Data[0][2])
//	recovered_stock.Price = value
//	return recovered_stock, true
//}

var wg sync.WaitGroup
var currentWaitGroup int64
var maxRoutines int64

func LoadAllStocks(timeInterval string) {
	//because the free version of alpha-vantage, has a limitation, also limit the concrrent routines fetching stocks
	maxRoutines = 4
	var currentWaitGroup int64
	var stocks []stock.Stock
	stocks = stock.LoadAllStockSymbols(timeInterval)
	logger.LogMessage("Starting to fetch stocks", logger.INFO)
	for index := range stocks {
		wg.Add(1)
		currentWaitGroup++
		LoadAndStoreStock(stocks[index])
		if currentWaitGroup >= maxRoutines {
			currentWaitGroup = 0
			time.Sleep(1 * time.Minute)
			wg.Wait()
		}
	}
	currentWaitGroup = 0
	wg.Wait()
	logger.LogMessage("Finished loading stocks", logger.INFO)
}

func LoadAndStoreStock(stocking stock.Stock) {
	defer wg.Done()
	stock, success := GetStockDataForSymbol(stocking, stock.ConvertTimeSeries(stocking.TimeData))
	fmt.Println(stock)
	if success {
		stock.Store()
	}
}
