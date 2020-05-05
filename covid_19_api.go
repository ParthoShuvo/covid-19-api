package main

import (
	"net/http"
	"time"

	"github.io/covid-19-api/middleware"
	"github.io/covid-19-api/model"
	"github.io/covid-19-api/resource"
	"github.io/covid-19-api/resource/writer"
	"github.io/covid-19-api/route"
	"github.io/covid-19-api/uc/country"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/cfg"
	"github.io/covid-19-api/log4u"
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
		log.Panic("DB initialization failed")
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

func buildRoute(db *model.DB) *mux.Router {
	rb := route.NewRouteBuilder(config.AllowCORS(), config.AppName(), log4u.ContainsLogDebug(config.Logging().Level))
	apirb := rb.NewSubrouteBuilder("/covid-19/api")
	addCountryRoutes(apirb, db)
	addCSSERoutes(apirb, db)
	return rb.Router()
}

func addCountryRoutes(rb *route.Builder, db *model.DB) {
	var emptyQry map[string]string
	cres := resource.NewCountryResource(country.NewEnv(db), writer.NewWriter(writer.JSON))
	rb.Add("Countries", []string{http.MethodGet}, "/countries", emptyQry, cres.CountryFetcher())
	rb.Add("CountryByCC", []string{http.MethodGet}, "/countries/cc/{cc:[a-zA-Z][a-zA-Z]}", emptyQry, cres.CountryFetcherByCC())
	rb.Add("CountryByName", []string{http.MethodGet}, "/countries/name/{name:[a-zA-Z]+}", emptyQry, cres.CountryFetcherByName())
}

func addCSSERoutes(rb *route.Builder, db *model.DB) {
	// var emptyQuery map[string]string
	// csseres := resource.NewCsseDailyReportsResource(db.NewDataAccessor(db.CsseDailyData), writer.NewWriter(writer.JSON))
	// csserb := rb.NewSubrouteBuilder("/csse")
	// csserb.Add("DailyReports", []string{http.MethodGet}, "/daily-reports", emptyQuery, csseres.DailyReportsFetcher())

	// datePattern := "[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]" //MM-dd-YYYY
	// datePath := fmt.Sprintf("{date:%s}", datePattern)
	// csserb.Add("DailyReportsByDate", []string{http.MethodGet}, "/daily-reports/"+datePath, emptyQuery, csseres.DailyReportsFetcherByDate())
}

func addMiddleware(router *mux.Router) http.Handler {
	plm := middleware.PerformanceLogMiddleware(router, log4u.ContainsLogDebug(config.Logging().Level))
	return middleware.CORSMiddleware(plm, config.AllowCORS(), config.CORS())
}
