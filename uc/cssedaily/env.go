package cssedaily

import "github.io/covid-19-api/uc"

type Env struct {
	ds CsseDailyDatastore
}

func New(ds CsseDailyDatastore) *Env {
	return &Env{ds}
}

const (
	defaultPageSize int = 10
	defaultPage     int = 1
)

func (env *Env) ReadAllDailyReports(page, pageSize int) ([]*DailyReport, error) {
	if dailyReports, err := env.ds.ReadAllDailyReports(); err != nil {
		return dailyReports, err
	} else {
		paginator := uc.NewPaginator(defaultPageSize, defaultPage)
		offset, limit := paginator.Paginate(len(dailyReports), pageSize, page)
		return dailyReports[offset:limit:limit], err
	}
}

func (env *Env) ReadDailyReports(date string) (*DailyReport, error) {
	return env.ds.ReadDailyReport(date)
}
