package main

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Liberatys/libra-back/main/stock"
	"github.com/gorilla/mux"
)

var stockData []stock.Stock
var serializedStockItems []byte

type StockHolderstruct struct {
	Stocks []stock.Stock `json:"stocks"`
}

func PurgeStockScreen() {
	stockData = nil
}

func GetStocks(w http.ResponseWriter, r *http.Request) {
	if len(stockData) > 0 && stockData != nil {
		w.Write(serializedStockItems)
		return
	}
	stocks := stock.LoadStocksForRoute("5")
	for key := range stocks {
		stocks[key].Load()
	}
	stockHolderInstance := StockHolderstruct{
		stocks,
	}
	data, err := json.Marshal(stockHolderInstance)
	if err != nil {
		fmt.Println(err.Error())
	}
	serializedStockItems = data
	stockData = stocks
	w.Write(serializedStockItems)
}

func GetStockByParameter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	symbol := params["symbol"]
	interval := params["interval"]
	stockInstance := stock.Stock{
		Symbol:   symbol,
		TimeData: interval,
	}
	stockInstance.Load()
	w.Write([]byte(stockInstance.Data))
}

func compress(source string) []byte {
	buf := new(bytes.Buffer)
	w, _ := flate.NewWriter(buf, 7)
	w.Write([]byte(source))
	w.Flush()
	return buf.Bytes()
}
