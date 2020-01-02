package apiconnection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/stock"
	av "github.com/cmckee-dev/go-alpha-vantage"
)

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

func LoadAllStocks(timeInterval string) {
	var wg sync.WaitGroup
	var currentWaitGroup int64
	var maxRoutines int64
	//because the free version of alpha-vantage, has a limitation, also limit the concrrent routines fetching stocks
	maxRoutines = 5
	var stocks []stock.Stock
	stocks = stock.LoadAllStockSymbols(timeInterval)
	logger.LogMessage("Starting to fetch stocks", logger.INFO)
	for index := range stocks {
		wg.Add(1)
		currentWaitGroup++
		LoadAndStoreStock(stocks[index], &wg)
		if currentWaitGroup >= maxRoutines {
			currentWaitGroup = 0
			time.Sleep(1 * time.Minute)
			wg.Wait()
		}
	}
	wg.Wait()
	currentWaitGroup = 0
	logger.LogMessage("Finished loading stocks", logger.INFO)
}

func LoadAndStoreStock(stocking stock.Stock, wg *sync.WaitGroup) {
	defer wg.Done()
	stock, success := GetStockDataForSymbol(stocking, stock.ConvertTimeSeries(stocking.TimeData))
	if success {
		logger.LogMessage(fmt.Sprintf("Stock %s was loaded", stock.Symbol), logger.INFO)
		if stock.Company == "" {
			stock.Company = GetCompanyNameForSymbol(stocking.Symbol)
		}
		stock.Update()
	}
}

type YahooReponse struct {
	ResultSet struct {
		Query  string `json:"Query"`
		Result []struct {
			Symbol   string `json:"symbol"`
			Name     string `json:"name"`
			Exch     string `json:"exch"`
			Type     string `json:"type"`
			ExchDisp string `json:"exchDisp"`
			TypeDisp string `json:"typeDisp"`
		} `json:"Result"`
	} `json:"ResultSet"`
}

func GetCompanyNameForSymbol(symbol string) string {
	url := fmt.Sprintf("http://d.yimg.com/autoc.finance.yahoo.com/autoc?query=%s&region=1&lang=en", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "-"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "-"
	}
	var yahooResponse YahooReponse
	err = json.Unmarshal(body, &yahooResponse)
	if err != nil {
		fmt.Println(err.Error())
		return "-"
	}
	return yahooResponse.ResultSet.Result[0].Name
}
