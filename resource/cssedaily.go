package resource

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/errors"
	"github.io/covid-19-api/resource/writer"
	"github.io/covid-19-api/uc/cssedaily"
)

type (
	// CsseDailyReportsResource defines country resources
	CsseDailyReportsResource struct {
		env    *cssedaily.Env
		writer writer.Writer
	}

	pageInfo struct {
		page     int
		pageSize int
	}
)

// NewCsseDailyReportsResource definition
func NewCsseDailyReportsResource(env *cssedaily.Env, w writer.Writer) *CsseDailyReportsResource {
	return &CsseDailyReportsResource{env, w}
}

// DailyReportsFetcher provides action to fetch all dailyreports
func (cdr *CsseDailyReportsResource) DailyReportsFetcher() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		pinfo, pageInfoErr := cdr.parsePageInfo(req.URL.Query())
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

func (cdr *CsseDailyReportsResource) parsePageInfo(query url.Values) (*pageInfo, error) {
	var (
		page            int
		pageSize        int
		defaultPage     int = 1
		defaultPageSize int = 10
	)
	anyOrElse := func(val string, defaultVal int) int {
		if parsedVal, err := strconv.Atoi(val); err == nil {
			return parsedVal
		}
		return defaultVal
	}
	isNeg := func(val int) bool { return val < 0 }
	page, pageSize = anyOrElse(query.Get("page"), defaultPage), anyOrElse(query.Get("pagesize"), defaultPageSize)
	if isNeg(page) || isNeg(pageSize) {
		return &pageInfo{}, errors.BadRequest.New("malformed query param values in page & pagesize")
	}
	return &pageInfo{page, pageSize}, nil
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
