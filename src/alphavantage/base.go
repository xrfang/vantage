package alphavantage

import (
	"net/http"
	"time"
)

type OHLCV struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type AVReply struct {
	Error    string `json:"Error Message"`
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		Interval      string `json:"4. Interval"`
		OutputSize    string `json:"5. Output Size"`
		TimeZone      string `json:"6. Time Zone"`
	} `json:"Meta Data"`
	TimeSeries_1min map[string]OHLCV `json:"Time Series (1min)"`
	TimeSeries_5min map[string]OHLCV `json:"Time Series (5min)"`
}

type client struct {
	hc  http.Client
	url string
	key string
}

func (c *client) SetUrl(url string) {
	c.url = url
}

func (c *client) SetTimeout(timeout int) {
	c.hc.Timeout = time.Duration(timeout) * time.Second
}

func NewClient(key string) client {
	return client{
		hc:  http.Client{Timeout: 60 * time.Second},
		url: "http://www.alphavantage.co",
		key: key,
	}
}
