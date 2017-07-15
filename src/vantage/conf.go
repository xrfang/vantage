package main

import (
	"github.com/xrfang/go-conf"
)

type Configuration struct {
	HTTP_PORT     string
	READ_TIMEOUT  int
	WRITE_TIMEOUT int
	API_KEY       string
}

var cf Configuration

func loadConfig(fn string) {
	//default values
	cf.HTTP_PORT = "8684"
	cf.READ_TIMEOUT = 60
	cf.WRITE_TIMEOUT = 60
	if fn != "" {
		assert(conf.ParseFile(fn, &cf))
	}
}
