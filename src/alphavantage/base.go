package alphavantage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OHLC struct {
	Open       string `json:"1. open"`
	High       string `json:"2. high"`
	Low        string `json:"3. low"`
	Close      string `json:"4. close"`
	Volume     string `json:"5. volume"`
	AdjClose   string `json:"5. adjusted close"`
	AdjVolume  string `json:"6. volume"`
	Dividend   string `json:"7. dividend amount"`
	SplitCoeff string `json:"8. split coefficient"`
}

type AVReply struct {
	Error    string `json:"Error Message"`
	MetaData struct {
		Information   string `json:"1. Information"`
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		TZ1           string `json:"4. Time Zone"`
		TZ2           string `json:"5. Time Zone"`
		TZ3           string `json:"6. Time Zone"`
	} `json:"Meta Data"`
	TimeSeries_1min    map[string]OHLC `json:"Time Series (1min)"`
	TimeSeries_5min    map[string]OHLC `json:"Time Series (5min)"`
	TimeSeries_15min   map[string]OHLC `json:"Time Series (15min)"`
	TimeSeries_30min   map[string]OHLC `json:"Time Series (30min)"`
	TimeSeries_60min   map[string]OHLC `json:"Time Series (60min)"`
	TimeSeries_Daily   map[string]OHLC `json:"Time Series (Daily)"`
	TimeSeries_Weekly  map[string]OHLC `json:"Weekly Time Series"`
	TimeSeries_Monthly map[string]OHLC `json:"Monthly Time Series"`
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
		url: "http://www.alphavantage.co/query",
		key: key,
	}
}

func (c client) request(args map[string]string) (*AVReply, error) {
	var params []string
	for k, v := range args {
		params = append(params, k+"="+v)
	}
	url := c.url + "?" + strings.Join(params, "&")
	fmt.Println(url)
	resp, err := c.hc.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}
	dec := json.NewDecoder(resp.Body)
	var ar AVReply
	err = dec.Decode(&ar)
	if err != nil {
		return nil, err
	}
	if ar.Error != "" {
		return nil, fmt.Errorf(ar.Error)
	}
	return &ar, nil
}
