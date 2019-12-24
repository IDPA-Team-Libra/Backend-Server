package main

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Liberatys/libra-back/main/stock"
)

var stock_data []stock.Stock
var stock_string []byte

type StockHolderstruct struct {
	Stocks []stock.Stock `json:"stocks"`
}

func GetStocks(w http.ResponseWriter, r *http.Request) {
	if len(stock_data) > 0 {
		w.Write(stock_string)
		return
	}
	stocks := stock.LoadStocksForRoute("5")
	for key := range stocks {
		stocks[key].Load()
	}
	stock_dat := StockHolderstruct{
		stocks,
	}
	data, err := json.Marshal(stock_dat)
	if err != nil {
		fmt.Println(err.Error())
	}
	stock_string = data
	stock_data = stocks
	w.Write(stock_string)
}

func compress(source string) []byte {
	buf := new(bytes.Buffer)
	w, _ := flate.NewWriter(buf, 7)
	w.Write([]byte(source))
	w.Flush()
	return buf.Bytes()
}
