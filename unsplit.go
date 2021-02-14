package main

import (
	"fmt"
	"strconv"
	"strings"
)

func unsplit(data []map[string]string) (adj []map[string]string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	s := func(v float64) string {
		sv := fmt.Sprintf("%f", v)
		sv = strings.TrimRight(sv, "0")
		return strings.TrimRight(sv, ".")
	}
	for _, d := range data {
		f := func(key string) float64 {
			v, err := strconv.ParseFloat(d[key], 64)
			assert(err)
			return v
		}
		r := f("adj_close") / f("close")
		d["open"] = s(f("open") * r)
		d["high"] = s(f("high") * r)
		d["low"] = s(f("low") * r)
		d["close"] = d["adj_close"]
		delete(d, "adj_close")
		delete(d, "split_coeff")
		delete(d, "dividend")
		adj = append(adj, d)
	}
	return
}
