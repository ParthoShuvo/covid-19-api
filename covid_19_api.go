package main

import (
	"fmt"
	"net/http"

	"github.io/covid-19-api/cfg"
)

func main() {
	config := cfg.NewConfig(version)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "api works!")
	})

	http.ListenAndServe(config.Server().String(), nil)
}
