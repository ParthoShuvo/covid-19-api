package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/cfg"
	"github.io/covid-19-api/log4u"
)

func main() {
	config := cfg.NewConfig(version)
	log4u.ConfigureLogging(config.Logging().Filename, config.Logging().Level)
	defer log4u.CloseLog()
	log.Infof("Starting %s on %s", config.AppName(), config.Server().String())
	log.Fatal(http.ListenAndServe(config.Server().String(), nil))
}
