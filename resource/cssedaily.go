package resource

import (
	"net/http"

	"github.io/covid-19-api/db"
	"github.io/covid-19-api/resource/writer"
)

// CsseDailyReportsResource defines country resources
type CsseDailyReportsResource struct {
	da     db.DataAccessor
	writer writer.Writer
}

// NewCsseDailyReportsResource definition
func NewCsseDailyReportsResource(dataAccessor db.DataAccessor, w writer.Writer) *CsseDailyReportsResource {
	return &CsseDailyReportsResource{dataAccessor, w}
}

// DailyReportsFetcher provides action to fetch all dailyreports
func (cdr *CsseDailyReportsResource) DailyReportsFetcher() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		dailyReportsMap := cdr.da.GetAll().(map[string]*db.CsseDailyReports)
		cdr.writer.Write(rw, dailyReportsMap)
	}
}
