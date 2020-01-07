package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Liberatys/libra-back/main/stock"
	"github.com/gorilla/mux"
)

var serializedStockItems []byte = nil

//Stocks is a serialized holder for stocks
type Stocks struct {
	Stocks string `json:"stocks"`
}

//PurgeStockScreen removes all memory stored stock values
func PurgeStockScreen() {
	serializedStockItems = nil
}

//GetStocks returns all stock items that are used for the market page
func GetStocks(w http.ResponseWriter, r *http.Request) {
	if serializedStockItems != nil {
		w.Write(serializedStockItems)
		return
	}
	stocks := stock.LoadStocksForRoute("5")
	for key := range stocks {
		stocks[key].Load()
		stocks[key].Data = ""
	}
	var b bytes.Buffer
	writer := zlib.NewWriter(&b)
	data, err := json.Marshal(stocks)
	writer.Write(data)
	writer.Close()
	encoded := base64.StdEncoding.EncodeToString([]byte(b.Bytes()))
	stockHolderInstance := Stocks{
		Stocks: encoded,
	}
	json, err := json.Marshal(stockHolderInstance)
	if err != nil {
		fmt.Println(err.Error())
	}
	serializedStockItems = json
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
