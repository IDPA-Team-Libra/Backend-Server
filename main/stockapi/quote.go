package stockapi

type MetaData struct {
	OneInformation     string `json:"1. Information"`
	TwoSymbol          string `json:"2. Symbol"`
	ThreeLastRefreshed string `json:"3. Last Refreshed"`
	FourInterval       string `json:"4. Interval"`
	FiveOutputSize     string `json:"5. Output Size"`
	SixTimeZone        string `json:"6. Time Zone"`
}

type QuoteResponse struct {
	MetaData   MetaData   `json:"Meta Data"`
	TimeSeries TimeSeries `json:"Time Series (5min)`
}
type TimeSeries struct {
}

type Instance struct {
	OneOpen    string `json:"1. open"`
	TwoHigh    string `json:"2. high"`
	ThreeLow   string `json:"3. low"`
	FourClose  string `json:"4. close"`
	FiveVolume string `json:"5. volume"`
}
