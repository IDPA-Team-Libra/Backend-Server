package apiconnection

import (
	"fmt"
	"sync"

	"github.com/DannyBen/quandl"
	"github.com/Liberatys/libra-back/main/stock"
	av "github.com/cmckee-dev/go-alpha-vantage"
	"github.com/robfig/cron"
)

//TODO hide api_key
const APIKEY = "F1HqA-xHQe7tzNYtFf26"

func GetStockDataForSymbol(recovered_stock stock.Stock, interval av.TimeInterval) (stock.Stock, bool) {
	quandl.APIKey = APIKEY
	data, err := quandl.GetSymbol("WIKI/"+recovered_stock.Symbol, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(data.Data) == 0 {
		return recovered_stock, false
	}
	value := fmt.Sprintf("%v", data.Data[0][2])
	recovered_stock.Price = value
	return recovered_stock, true
}

var wg sync.WaitGroup
var current_wait_group int64
var max_routines int64

// TODO check this code
// TODO check if it only executes given max_routines at a time
func LoadAllStocks(timeInterval string) {
	max_routines = 1
	var current_wait_group int64
	var stocks []stock.Stock
	stocks = stock.LoadAllStockSymbols(timeInterval)
	for _, value := range stocks {
		fmt.Println(value)
		wg.Add(1)
		current_wait_group += 1
		LoadAndStoreStock(value)
		if current_wait_group == max_routines {
			current_wait_group = 0
			wg.Wait()
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
