package resource

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/resource/writer"
	"github.io/covid-19-api/uc/cssedaily"
)

// CsseDailyReportsResource defines country resources
type CsseDailyReportsResource struct {
	env    *cssedaily.Env
	writer writer.Writer
}

// NewCsseDailyReportsResource definition
func NewCsseDailyReportsResource(env *cssedaily.Env, w writer.Writer) *CsseDailyReportsResource {
	return &CsseDailyReportsResource{env, w}
}

// DailyReportsFetcher provides action to fetch all dailyreports
func (cdr *CsseDailyReportsResource) DailyReportsFetcher() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		pinfo, pageInfoErr := parsePageInfo(req.URL.Query())
		if pageInfoErr != nil {
			log.Error(pageInfoErr)
			SendError(rw, pageInfoErr)
			return
		}
		dailyReports, err := cdr.env.ReadAllDailyReports(pinfo.page, pinfo.pageSize)
		if err != nil {
			log.Error(err)
			SendError(rw, err)
			return
		}
		cdr.writer.Write(rw, dailyReports)
	}
}

// DailyReportsFetcherByDate provides action to fetch one dailyreports by date
func (cdr *CsseDailyReportsResource) DailyReportsFetcherByDate() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		date := mux.Vars(req)["date"]
		dailyReports, err := cdr.env.ReadDailyReports(date)
		if err != nil {
			log.Error(err)
			SendError(rw, err)
			return
		}
		cdr.writer.Write(rw, dailyReports)
	}
}
