package stock

type PortfolioStock struct {
	ID                string  `json:"id"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	BuyPrice          float64 `json:"buyPrice"`
	BuyDate           string  `json:"buyDate"`
	StockID           int     `json:"stockID"`
	Amount            int64   `json:"amount "`
	MarketValue       float64 `json:"marketValue"`
	DifferencePercent float64 `json:"DifferencePercent"`
	Difference        float64 `json:"difference"`
	Result            bool    `json:"result"`
}

func NewStockPortfolioPrice(symbol string) PortfolioStock {
	portStock := PortfolioStock{
		Symbol: symbol,
	}
	return portStock
}

func (portStock *PortfolioStock) Store() {

}

func (portStock *PortfolioStock) Overwrite() {

}

func (portStock *PortfolioStock) Retreave() {

}
