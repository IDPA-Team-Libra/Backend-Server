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

func GetStockDataForSymbol(recovered_stock stock.Stock, key string) (stock.Stock, bool) {
	client := av.NewClient(key)
	interval, _ := stock.ConvertTimeSeries(recovered_stock.TimeData)
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
	json, _ := json.Marshal(result)
	recovered_stock.Price = price
	recovered_stock.Data = string(json)
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
	keyRounds = []string{
		"CG96DXD2YPARDLMX",
		"GJQQ5PRFZSPZSENI",
		"9PW8B5GO7RHT3XJ5",
		"7JJ4KU81J5CDN1CU",
		"6CXB51K9FI74XVK",
		"PAKC21FC1QJR6YMX",
	}
	nameCache = make(map[string]string)
	//because the free version of alpha-vantage, has a limitation, also limit the concrrent routines fetching stocks
	maxRoutines = 30
	var stocks []stock.Stock
	timeIntervals := []string{
		"1",
		"5",
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
			if currentWaitGroup%5 == 0 {
				keyIndex++
				if keyIndex == len(keyRounds)-1 {
					keyIndex = 0
				}
			}
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

func LoadAndStoreStock(stocking stock.Stock, wg *sync.WaitGroup, key string) {
	defer wg.Done()
	stock, success := GetStockDataForSymbol(stocking, key)
	fmt.Println(success)
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
