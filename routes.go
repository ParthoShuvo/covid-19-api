package main

import (
	"fmt"
	"net/http"

	"github.com/ParthoShuvo/covid-19-api/log4u"
	"github.com/ParthoShuvo/covid-19-api/middleware"
	"github.com/ParthoShuvo/covid-19-api/model"
	"github.com/ParthoShuvo/covid-19-api/resource"
	"github.com/ParthoShuvo/covid-19-api/resource/writer"
	"github.com/ParthoShuvo/covid-19-api/route"
	"github.com/ParthoShuvo/covid-19-api/uc/country"
	"github.com/ParthoShuvo/covid-19-api/uc/cssedaily"
	"github.com/ParthoShuvo/covid-19-api/uc/cssetimeseries"
	"github.com/gorilla/mux"
)

var homeRouteHandler http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, config.HomeURL(), http.StatusSeeOther)
}

var apiRouteHandler http.HandlerFunc = func(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, config.APIDocURL(), http.StatusSeeOther)
}

func buildRoute(db *model.DB) *mux.Router {
	var emptyQry map[string]string
	rb := route.NewRouteBuilder(config.AllowCORS(), config.AppName(), log4u.ContainsLogDebug(config.Logging().Level))
	rb.Add("Home", []string{http.MethodGet, http.MethodHead}, "/", emptyQry, homeRouteHandler)
	apirb := rb.NewSubrouteBuilder("/api")
	apirb.Add("APIDoc", []string{http.MethodGet}, "", emptyQry, apiRouteHandler)
	addCountryRoutes(apirb, db)
	addCSSERoutes(apirb, db)
	return rb.Router()
}

func addCountryRoutes(rb *route.Builder, db *model.DB) {
	var emptyQry map[string]string
	cres := resource.NewCountryResource(country.NewEnv(db), writer.NewWriter(writer.JSON))
	rb.Add("CountriesByNames", []string{http.MethodGet}, "/countries", map[string]string{"name": "{[a-zA-Z ]*}"}, cres.CountryFetcherByName())
	rb.Add("Countries", []string{http.MethodGet}, "/countries", emptyQry, cres.CountryFetcher())
	rb.Add("CountryByCC", []string{http.MethodGet}, "/countries/cc/{cc:[a-zA-Z][a-zA-Z]}", emptyQry, cres.CountryFetcherByCC())
}

func addCSSERoutes(rb *route.Builder, db *model.DB) {
	var emptyQuery map[string]string
	res := resource.NewCsseDailyReportsResource(cssedaily.New(db), writer.NewWriter(writer.JSON))
	csserb := rb.NewSubrouteBuilder("/csse")
	paginationQueryParams := map[string]string{"page": "{[1-9][0-9]*}", "pagesize": "{[1-9][0-9]*}"}
	csserb.Add("DailyReportsWithPagination", []string{http.MethodGet}, "/daily-reports", paginationQueryParams, res.DailyReportsFetcher())

	csserb.Add("DailyReportsWithPagination", []string{http.MethodGet}, "/daily-reports", emptyQuery, res.DailyReportsFetcher())

	datePattern := "[0-9][0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]" //MM-dd-YYYY
	datePath := fmt.Sprintf("{date:%s}", datePattern)
	csserb.Add("DailyReportsByDate", []string{http.MethodGet}, "/daily-reports/"+datePath, emptyQuery, res.DailyReportsFetcherByDate())

	addCsseTimeSeriesRoutes(csserb, db)
}

func addCsseTimeSeriesRoutes(rb *route.Builder, db *model.DB) {
	res := resource.NewCsseTimeSeriesResource(cssetimeseries.New(db), writer.NewWriter(writer.JSON))
	trb := rb.NewSubrouteBuilder("/time-series")
	datePattern := "[0-9]-[0-9][0-9]-[0-9][0-9][0-9][0-9]"

	paginateQueryParams := map[string]string{
		"start":    fmt.Sprintf("{%s}", datePattern), //MM-dd-YYYY
		"end":      "{\\d\\d-\\d\\d-\\d\\d\\d\\d}",   //MM-dd-YYYY
		"page":     "{[1-9][0-9]*}",
		"pagesize": "{[1-9][0-9]*}",
	}
	trb.Add("TimeSeriesInBetweenByPage", []string{http.MethodGet}, "/between", paginateQueryParams, res.TimeSeriesFetcherInBetween())

	paginateByPageQueryParams := map[string]string{
		"start": fmt.Sprintf("{%s}", datePattern), //MM-dd-YYYY
		"end":   "{\\d\\d-\\d\\d-\\d\\d\\d\\d}",   //MM-dd-YYYY
		"page":  "{[1-9][0-9]*}",
	}
	trb.Add("TimeSeriesInBetweenByPage", []string{http.MethodGet}, "/between", paginateByPageQueryParams, res.TimeSeriesFetcherInBetween())

	inBetwQueryParams := map[string]string{
		"start": fmt.Sprintf("{%s}", datePattern), //MM-dd-YYYY
		"end":   "{\\d\\d-\\d\\d-\\d\\d\\d\\d}",   //MM-dd-YYYY
	}
	trb.Add("TimeSeriesInBetween", []string{http.MethodGet}, "/between", inBetwQueryParams, res.TimeSeriesFetcherInBetween())

	timeSeriesQuery := map[string]string{"date": fmt.Sprintf("{%s}", datePattern)}
	trb.Add("TimeSeriesByDate", []string{http.MethodGet}, "/select", timeSeriesQuery, res.TimeSeriesFetcherByDate())

}

func addMiddleware(router *mux.Router) http.Handler {
	rm := middleware.RecoveryMiddleware(router)
	plm := middleware.PerformanceLogMiddleware(rm, log4u.ContainsLogDebug(config.Logging().Level))
	crsm := middleware.CORSMiddleware(plm, config.AllowCORS(), config.CORS())
	return crsm
}
