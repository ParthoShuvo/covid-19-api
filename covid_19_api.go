package main

import (
	"net/http"
	"time"

	"github.com/ParthoShuvo/covid-19-api/model"

	"github.com/ParthoShuvo/covid-19-api/cfg"
	"github.com/ParthoShuvo/covid-19-api/log4u"
	log "github.com/sirupsen/logrus"
)

var config *cfg.Config

func init() {
	config = cfg.NewConfig(version)
	log4u.ConfigureLogging(config.Logging().Filename, config.Logging().Level)
}

func main() {
	defer log4u.CloseLog()
	var (
		db  *model.DB
		err error
	)
	if db, err = model.NewDB(config.Database()); err != nil {
		log.Fatal("DB initialization failed")
	}
	srv := &http.Server{
		Addr:         config.Server().String(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      addMiddleware(buildRoute(db)),
	}

	log.Infof("Starting %s on %s", config.AppName(), config.Server().String())
	log.Fatal(srv.ListenAndServe())
}
