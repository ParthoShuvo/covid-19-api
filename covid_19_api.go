package main

import (
	"fmt"

	"github.io/covid-19-api/cfg"
)

func main() {
	config := cfg.NewConfig(version)
	fmt.Println(config.AppName(), "starts")
}
