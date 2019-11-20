package main

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Liberatys/libra-back/main/stock"
)

//TODO: Implement the methods required to store the information in angular
var stock_data []stock.Stock
var stock_string []string

func GetStocks(w http.ResponseWriter, r *http.Request) {
	if len(stock_data) > 0 {
		w.Write([]byte(stock_string))
	}
	stocks := stock.LoadAllStockSymbols("5")
	for key := range stocks {
		stocks[key].Load()
		data, err := json.Marshal(stocks[key])
		if err != nil {

		}
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	stock_data = stocks
}

func compress(source string) []byte {
	buf := new(bytes.Buffer)
	w, _ := flate.NewWriter(buf, 7)
	w.Write([]byte(source))
	w.Flush()
	return buf.Bytes()
}
