package apiconnection

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Liberatys/libra-back/main/stock"
	av "github.com/cmckee-dev/go-alpha-vantage"
	"github.com/robfig/cron"
)

//TODO hide api_key
const APIKEY = "CG96DXD2YPARDLMX"

func GetStockDataForSymbol(recovered_stock stock.Stock, interval av.TimeInterval) (stock.Stock, bool) {
	client := av.NewClient(APIKEY)
	result, err := client.StockTimeSeriesIntraday(interval, recovered_stock.Symbol)
	if err != nil {
		return stock.Stock{}, false
	}
	if len(result) == 0 {
		return recovered_stock, false
	}
	price := fmt.Sprintf("%.3f", result[len(result)-1].Close)
	pagesJson, err := json.Marshal(result)
	if err != nil {
		return stock.Stock{}, false
	}
	recovered_stock.Price = price
	recovered_stock.Data = string(pagesJson)
	return recovered_stock, true
}

var wg sync.WaitGroup
var current_wait_group int64
var max_routines int64

// TODO check this code
// TODO check if it only executes given max_routines at a time
func LoadAllStocks(timeInterval string) {
	max_routines = 5
	var stocks []stock.Stock
	stocks = stock.LoadAllStockSymbols(timeInterval)
	max_id := len(stocks) - 1
	for id, value := range stocks {
		wg.Add(1)
		current_wait_group += 1
		go LoadAndStoreStock(value)
		if max_id-id > 10 {
			if current_wait_group >= max_routines {
				current_wait_group = 0
				wg.Wait()
			}
		}
	}
	current_wait_group = 0
	wg.Wait()
}

func LoadAndStoreStock(stocking stock.Stock) {
	defer wg.Done()
	stock, success := GetStockDataForSymbol(stocking, stock.ConvertTimeSeries(stocking.TimeData))
	if success {
		stock.Store()
	}
}

// TODO validate the execution of the cron jobs
func StartCronJobs() {
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		LoadAllStocks("5")
	})
	c.Start()
}
