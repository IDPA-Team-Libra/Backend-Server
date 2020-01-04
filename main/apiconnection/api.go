package apiconnection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	av "github.com/Liberatys/go-alpha-vantage"
	"github.com/Liberatys/libra-back/main/logger"
	"github.com/Liberatys/libra-back/main/stock"
)

var keyRounds []string

//GetStockDataForSymbol loads the stock price and data from the alpha-vantage api
func GetStockDataForSymbol(recoveredStock stock.Stock, key string) (stock.Stock, bool) {
	client := av.NewClient(key)
	interval, _ := stock.ConvertTimeSeries(recoveredStock.TimeData)
	result, err := client.StockTimeSeriesIntraday(interval, recoveredStock.Symbol)
	if err != nil {
	}
	if err != nil {
		return stock.Stock{}, false
	}
	if len(result) == 0 {
		return recoveredStock, false
	}
	price := fmt.Sprintf("%.3f", result[len(result)-1].Close)
	json, _ := json.Marshal(result)
	recoveredStock.Price = price
	recoveredStock.Data = string(json)
	return recoveredStock, true
}

//LoadAllStocks loads the stock information for all stocks in the database
func LoadAllStocks(timeInterval string) {
	var wg sync.WaitGroup
	var currentWaitGroup int64
	var maxRoutines int64
	keyRounds = []string{
		"CG96DXD2YPARDLMX",
	}
	nameCache = make(map[string]string)
	//because the free version of alpha-vantage, has a limitation, also limit the concrrent routines fetching stocks
	maxRoutines = 5
	var stocks []stock.Stock
	timeIntervals := []string{
		"5",
		"1",
		"15",
		"30",
		"60",
	}
	var keyIndex int
	keyIndex = 0
	for index := range timeIntervals {
		stocks = stock.LoadAllStockSymbols(timeIntervals[index])
		logger.LogMessage("Starting to fetch stocks", logger.INFO)
		for index := range stocks {
			wg.Add(1)
			currentWaitGroup++
			LoadAndStoreStock(stocks[index], &wg, keyRounds[keyIndex])
			if currentWaitGroup >= maxRoutines {
				currentWaitGroup = 0
				time.Sleep(1 * time.Minute)
				wg.Wait()
			}
		}
	}
	wg.Wait()
	currentWaitGroup = 0
	//get the memory back
	nameCache = make(map[string]string)
	logger.LogMessage("Finished loading stocks", logger.INFO)
}

var nameCache map[string]string

//LoadAndStoreStock load stock information and write the new data to the database
func LoadAndStoreStock(stocking stock.Stock, wg *sync.WaitGroup, key string) {
	defer wg.Done()
	stock, success := GetStockDataForSymbol(stocking, key)
	if success {
		logger.LogMessage(fmt.Sprintf("Stock %s was loaded", stock.Symbol), logger.INFO)
		if stock.Company == "" {
			value, ok := nameCache[stock.Symbol]
			if ok == true {
				stock.Company = value
			} else {
				stock.Company = GetCompanyNameForSymbol(stocking.Symbol)
				nameCache[stocking.Symbol] = stock.Company
			}
		}
		stock.Update()
	}
}

//YahooReponse a struct for parsing the yahoo response to resolve company names
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

//GetCompanyNameForSymbol returns the company name for a given stock-symbol -- only executed if stock has no company set [will only execute the first iteration]
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
