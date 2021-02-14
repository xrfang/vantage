package main

import (
	"bytes"
	"encoding/csv"
	"strings"
)

func tabulate(data []map[string]string) (string, string) {
	if len(data) == 0 {
		return "", ""
	}
	keys := []string{"time", "open", "high", "low", "close", "volume"}
	tab := []string{`<table border="0" width="100%"><tr class="headrow">`}
	for _, k := range keys {
		tab = append(tab, `<th class="thcell">`+k+`</th>`)
	}
	tab = append(tab, `</tr>`)
	n := 0
	for i := len(data) - 1; i >= 0; i-- {
		d := data[i]
		if i%2 == 0 {
			tab = append(tab, `<tr class="evenrow">`)
		} else {
			tab = append(tab, `<tr class="oddrow">`)
		}
		for _, k := range keys {
			tab = append(tab, `<td class="tdcell">`+d[k]+`</td>`)
		}
		tab = append(tab, `</tr>`)
		n++
		if n >= 100 {
			break
		}
	}
	var buf bytes.Buffer
	cw := csv.NewWriter(&buf)
	cw.UseCRLF = true
	cw.Write(keys)
	for _, d := range data {
		var row []string
		for _, k := range keys {
			row = append(row, d[k])
		}
		cw.Write(row)
	}
	cw.Flush()
	return strings.Join(tab, "\n"), buf.String()
}
