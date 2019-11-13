package stockapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
)

func SendRequest() {
	resp, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=MSFT&interval=5min&apikey=demo")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(jsonParsed)
}
