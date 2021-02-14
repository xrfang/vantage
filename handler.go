package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"vantage/alphavantage"

	"github.com/atotto/clipboard"
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
<form method=GET onsubmit="doQuery()">
<div style="width:100%;position:absolute;padding:6px;background:#666666">
<span style="color:white">
SYMBOL&nbsp;
<input style="height:32px" name="symbol" type="text" value={{SYMBOL}}>
<input style="height:32px" id="qry" type="submit" />
</span>
<input name=compact id=compact type="checkbox" {{COMPACT}} />
<a href="javascript:switchSize()" style="color:white;text-decoration:none"
>Latest 100 data points only</a>
<span style="float:right;margin-right:16px;color:white">PERIOD&nbsp;
<select name="period" style="height:32px">
<option value="0" {{SEL0}}>1 min</opton>
<option value="1" {{SEL1}}>5 min</opton>
<option value="2" {{SEL2}}>15 min</opton>
<option value="3" {{SEL3}}>30 min</opton>
<option value="4" {{SEL4}}>60 min</opton>
<option value="5" {{SEL5}}>Daily</opton>
<option value="6" {{SEL6}}>Daily Adjusted</opton>
<option value="7" {{SEL7}}>Weekly</opton>
<option value="8" {{SEL8}}>Monthly</opton>
</select></span>
</div>
</form>
<div style="position:absolute;margin-top:44px;width:100%;height:100%">
<div style="background:lightyellow;padding:6px;font-weight:bold">{{SUMMARY}}</div>
{{CONTENT}}
</div>
<script>
function doQuery() {
    document.getElementById("qry").disabled = true;
}
function switchSize() {
	compact = document.getElementById("compact")
	compact.checked = !compact.checked
}
</script>
</body>
</html>
`
)

func handler(w http.ResponseWriter, r *http.Request) {
	var content, summary string
	compact := "checked"
	output := "compact"
	symbol := strings.ToUpper(r.URL.Query().Get("symbol"))
	period := r.URL.Query().Get("period")
	p, err := strconv.Atoi(period)
	sel := "{{SEL6}}"
	if err == nil {
		sel = fmt.Sprintf("{{SEL%d}}", p)
	}
	if symbol != "" {
		if r.URL.Query().Get("compact") == "" {
			compact = ""
			output = "full"
		}
		v := alphavantage.NewClient(cf.API_KEY)
		data, meta, err := v.TimeSeries(alphavantage.TimeSeriesKind(p),
			symbol, output)
		if err != nil {
			summary = err.Error()
		} else {
			tpl := `<span><span style="color:red">%s</span>&nbsp;%s (%d 
			    items)</span><span style="float:right">%s (%s)</span>`
			caption := meta["symbol"]
			info := meta["information"]
			table, csv := tabulate(data)
			content = table
			err := clipboard.WriteAll(symbol + "\n" + csv)
			if err != nil {
				caption = "ERROR"
				info = "clipboard access failed: " + err.Error()
			}
			summary = fmt.Sprintf(tpl, caption, info, len(data),
				meta["last_refreshed"], meta["time_zone"])
		}
	}
	page := strings.Replace(PAGE, "{{VERSION}}", fmt.Sprintf("V%s.%s",
		_G_REVS, _G_HASH), 1)
	page = strings.Replace(page, "{{CONTENT}}", content, 1)
	page = strings.Replace(page, "{{SYMBOL}}", symbol, 1)
	page = strings.Replace(page, "{{COMPACT}}", compact, 1)
	page = strings.Replace(page, "{{SUMMARY}}", summary, 1)
	for i := 0; i < 9; i++ {
		key := fmt.Sprintf("{{SEL%d}}", i)
		val := ""
		if key == sel {
			val = "selected"
		}
		page = strings.Replace(page, key, val, 1)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, page)
}
