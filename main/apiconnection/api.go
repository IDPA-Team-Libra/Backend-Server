package apiconnection

import (
	"fmt"

	av "github.com/cmckee-dev/go-alpha-vantage"
)

const APIKEY = "CG96DXD2YPARDLMX"

func GetStockDataForSymbol(symbol string, interval av.TimeInterval) {
	client := av.NewClient(APIKEY)
	result, err := client.StockTimeSeriesIntraday(interval, symbol)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, val := range result {
		fmt.Println(val)
	}
}
