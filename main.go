package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	ver := flag.Bool("version", false, "show version info")
	conf := flag.String("conf", "", "configuration file")
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	loadConfig(*conf)

	http.DefaultServeMux.HandleFunc("/", handler)
	svr := http.Server{
		Addr:         ":" + cf.HTTP_PORT,
		ReadTimeout:  time.Duration(cf.READ_TIMEOUT) * time.Second,
		WriteTimeout: time.Duration(cf.WRITE_TIMEOUT) * time.Second,
	}
	assert(svr.ListenAndServe())
}
