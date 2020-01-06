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

var serializedStockItems []byte = nil

//Stocks is a serialized holder for stocks
type Stocks struct {
	Stocks []stock.Stock `json:"stocks"`
}

//PurgeStockScreen removes all memory stored stock values
func PurgeStockScreen() {
	serializedStockItems = nil
}

//GetStocks returns all stock items that are used for the market page
func GetStocks(w http.ResponseWriter, r *http.Request) {
	if serializedStockItems != nil {
		w.Write(serializedStockItems)
	}
	stocks := stock.LoadStocksForRoute("5")
	for key := range stocks {
		stocks[key].Load()
	}
	stockHolderInstance := Stocks{
		stocks,
	}
	data, err := json.Marshal(stockHolderInstance)
	if err != nil {
		fmt.Println(err.Error())
	}
	serializedStockItems = data
	w.Write(serializedStockItems)
}

//GetStockByParameter returns information for the grahpics backend in dash
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

// method that comes to use in a later iteration, where more stocks are transfered
func compress(source string) []byte {
	buf := new(bytes.Buffer)
	w, _ := flate.NewWriter(buf, 7)
	w.Write([]byte(source))
	w.Flush()
	return buf.Bytes()
}
