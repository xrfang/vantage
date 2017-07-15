package main

import (
	"alphavantage"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	PAGE = `<html>
<head>
<title>Alpha Vantage {{VERSION}}</title>
<style>
.headrow {background:#666666;color:white}
.evenrow {background:#f8f8f8}
.oddrow {background:#e8e8e8}
.thcell {padding:6px;text-align:left}
.tdcell {padding:6px;vertical-align:top}
</style>
</head>
<body style="margin:0">
<form method=GET>
<div style="width:100%;position:absolute;padding:6px;background:#666666">
<span style="color:white">
SYMBOL&nbsp;<input style="height:32px" name="symbol" type="text">
<input style="height:32px" name="submit" type="submit" />
</span>
<input name=output type="checkbox" checked/><span style="color:white">Latest data only</span>
<span style="float:right;margin-right:16px;color:white">PERIOD&nbsp;
<select name="period" style="height:32px">
<option value="0">1 min</opton>
<option value="1">5 min</opton>
<option value="2">15 min</opton>
<option value="3">30 min</opton>
<option value="4">60 min</opton>
<option value="5">Daily</opton>
<option value="6" selected>Daily Adjusted</opton>
<option value="7">Weekly</opton>
<option value="8">Monthly</opton>
</select></span>
</div>
</form>
<div style="position:absolute;margin-top:44px;width:100%;height:100%">
{{CONTENT}}
</div>
</body>
</html>
`
)

func handler(w http.ResponseWriter, r *http.Request) {
	page := strings.Replace(PAGE, "{{VERSION}}", fmt.Sprintf("V%s.%s",
		_G_REVS, _G_HASH), 1)
	content := ""
	if r.URL.Query().Get("submit") != "" {
		symbol := r.URL.Query().Get("symbol")
		period, _ := strconv.Atoi(r.URL.Query().Get("period"))
		v := alphavantage.NewClient(cf.API_KEY)
		data, meta, err := v.TimeSeries(alphavantage.TimeSeriesKind(period),
			symbol, "compact")
		if err != nil {
			content = err.Error()
		} else {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			assert(enc.Encode(meta))
			assert(enc.Encode(data))
			content = buf.String()
		}
	}
	page = strings.Replace(page, "{{CONTENT}}", content, -1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, page)
	/*
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		data, err := c.TimeSeries(alphavantage.TS_Dy, "TSLAs", "compact")
		if err != nil {
			fmt.Fprintf(w, "ERROR: %v", err)
			return
		}
		enc := json.NewEncoder(w)
		err = enc.Encode(data)
		if err != nil {
			fmt.Fprintf(w, "ERROR: %v", err)
		}
	*/
}
