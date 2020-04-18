package main

import (
	"net/http"
	"time"

	"github.io/covid-19-api/middleware"
	"github.io/covid-19-api/route"

	"github.io/covid-19-api/resource/writer"

	"github.io/covid-19-api/resource"

	"github.io/covid-19-api/db"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/cfg"
	"github.io/covid-19-api/log4u"
)

var config *cfg.Config

func init() {
	config = cfg.NewConfig(version)
	log4u.ConfigureLogging(config.Logging().Filename, config.Logging().Level)
	db.AddConfig(config.Database())
}

func main() {
	defer log4u.CloseLog()

	srv := &http.Server{
		Addr:         config.Server().String(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      addMiddleware(buildRoute()),
	}

	log.Infof("Starting %s on %s", config.AppName(), config.Server().String())
	log.Fatal(srv.ListenAndServe())
}

func buildRoute() *mux.Router {
	rb := route.NewRouteBuilder(config.AllowCORS(), config.AppName(), log4u.ContainsLogDebug(config.Logging().Level))

	apirb := rb.NewSubrouteBuilder("/covid-19/api")

	cres := resource.NewCountryResource(db.NewDataAccessor(db.CountryData), writer.NewWriter(writer.JSON))
	apirb.Add("Countries", []string{http.MethodGet}, "/countries", cres.CountryFetcher())
	apirb.Add("CountryByCC", []string{http.MethodGet}, "/countries/cc/{cc:[a-zA-Z][a-zA-Z]}", cres.CountryFetcherByCC())
	apirb.Add("CountryByName", []string{http.MethodGet}, "/countries/name/{name:[a-zA-Z]+}", cres.CountryFetcherByName())

	csseres := resource.NewCsseDailyReportsResource(db.NewDataAccessor(db.CsseDailyData), writer.NewWriter(writer.JSON))
	csserb := apirb.NewSubrouteBuilder("/csse")
	csserb.Add("DailyReports", []string{http.MethodGet}, "/daily-reports", csseres.DailyReportsFetcher())
	csserb.Add("DailyReportsByDate", []string{http.MethodGet}, "/daily-reports/{date:[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]}", csseres.DailyReportsFetcherByDate())
	return rb.Router()
}

func addMiddleware(router *mux.Router) http.Handler {
	plm := middleware.PerformanceLogMiddleware(router, log4u.ContainsLogDebug(config.Logging().Level))
	return middleware.CORSMiddleware(plm, config.AllowCORS(), config.CORS())
}
