package resource

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"
// 	log "github.com/sirupsen/logrus"

// 	"github.io/covid-19-api/db"
// 	"github.io/covid-19-api/resource/writer"
// )

// // CsseDailyReportsResource defines country resources
// type CsseDailyReportsResource struct {
// 	da     db.DataAccessor
// 	writer writer.Writer
// }

// // NewCsseDailyReportsResource definition
// func NewCsseDailyReportsResource(dataAccessor db.DataAccessor, w writer.Writer) *CsseDailyReportsResource {
// 	return &CsseDailyReportsResource{dataAccessor, w}
// }

// // DailyReportsFetcher provides action to fetch all dailyreports
// func (cdr *CsseDailyReportsResource) DailyReportsFetcher() http.HandlerFunc {
// 	return func(rw http.ResponseWriter, req *http.Request) {
// 		dailyReportsMap := cdr.da.GetAll().(map[string]*db.CsseDailyReports)
// 		cdr.writer.Write(rw, dailyReportsMap)
// 	}
// }

// // DailyReportsFetcherByDate provides action to fetch one dailyreports by date
// func (cdr *CsseDailyReportsResource) DailyReportsFetcherByDate() http.HandlerFunc {
// 	return func(rw http.ResponseWriter, req *http.Request) {
// 		date := mux.Vars(req)["date"]
// 		i, err := cdr.da.GetOne(date)
// 		if err != nil {
// 			log.Errorf("no country found by name=%s", date)
// 			rw.WriteHeader(http.StatusNotFound)
// 			fmt.Fprint(rw, "Not found")
// 			return
// 		}
// 		dailyReports := i.(*db.CsseDailyReports)
// 		cdr.writer.Write(rw, map[string]*db.CsseDailyReports{dailyReports.Date: dailyReports})
// 	}
// }
