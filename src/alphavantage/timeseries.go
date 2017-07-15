package alphavantage

import (
	"sort"
)

type TimeSeriesKind int

const (
	TS_Me1 TimeSeriesKind = iota
	TS_Me5
	TS_Me15
	TS_Me30
	TS_Me60
	TS_Dy
	TS_DyAdj
	TS_Wk
	TS_Mh
)

func (c client) TimeSeries(kind TimeSeriesKind, symbol, output string) (
	[]map[string]string, map[string]string, error) {
	args := map[string]string{
		"function":   "TIME_SERIES_INTRADAY",
		"symbol":     symbol,
		"outputsize": output,
		"datatype":   "json",
		"apikey":     c.key,
	}
	switch kind {
	case TS_Me1:
		args["interval"] = "1min"
	case TS_Me5:
		args["interval"] = "5min"
	case TS_Me15:
		args["interval"] = "15min"
	case TS_Me30:
		args["interval"] = "30min"
	case TS_Me60:
		args["interval"] = "60min"
	case TS_Dy:
		args["function"] = "TIME_SERIES_DAILY"
	case TS_DyAdj:
		args["function"] = "TIME_SERIES_DAILY_ADJUSTED"
	case TS_Wk:
		args["function"] = "TIME_SERIES_WEEKLY"
	case TS_Mh:
		args["function"] = "TIME_SERIES_MONTHLY"
	}
	reply, err := c.request(args)
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]string{
		"information":    reply.MetaData.Information,
		"symbol":         reply.MetaData.Symbol,
		"last_refreshed": reply.MetaData.LastRefreshed,
		"interval":       reply.MetaData.Interval,
		"output_size":    reply.MetaData.OutputSize,
		"time_zone":      reply.MetaData.TimeZone,
	}
	var data []map[string]string
	extract := func(series map[string]OHLC) []map[string]string {
		var arr []map[string]string
		for k, v := range series {
			arr = append(arr, map[string]string{
				"time":        k,
				"open":        v.Open,
				"high":        v.High,
				"low":         v.Low,
				"close":       v.Close,
				"volume":      v.Volume + v.AdjVolume, //only one of them is non-empty
				"adj_close":   v.AdjClose,
				"dividend":    v.Dividend,
				"split_coeff": v.SplitCoeff,
			})
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i]["time"] < arr[j]["time"]
		})
		return arr
	}
	switch kind {
	case TS_Me1:
		data = extract(reply.TimeSeries_1min)
	case TS_Me5:
		data = extract(reply.TimeSeries_5min)
	case TS_Me15:
		data = extract(reply.TimeSeries_15min)
	case TS_Me30:
		data = extract(reply.TimeSeries_30min)
	case TS_Me60:
		data = extract(reply.TimeSeries_60min)
	case TS_Dy:
		data = extract(reply.TimeSeries_Daily)
	case TS_DyAdj:
		data = extract(reply.TimeSeries_Daily)
	case TS_Wk:
		data = extract(reply.TimeSeries_Weekly)
	case TS_Mh:
		data = extract(reply.TimeSeries_Monthly)
	}
	return data, meta, err
}
