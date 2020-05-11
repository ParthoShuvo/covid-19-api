package cssedaily

import "github.io/covid-19-api/uc"

type Env struct {
	ds CsseDailyDatastore
}

func New(ds CsseDailyDatastore) *Env {
	return &Env{ds}
}

const (
	defaultPageSize = 10
	defaultPage     = 1
)

func (env *Env) ReadAllDailyReports(page, pageSize int) (*CsseDailyData, error) {
	if dailyReports, err := env.ds.ReadAllDailyReports(); err != nil {
		return &CsseDailyData{}, err
	} else {
		paginator := uc.NewPaginator(defaultPageSize, defaultPage)
		totalSize := len(dailyReports)
		offset, limit := paginator.Paginate(totalSize, pageSize, page)
		return &CsseDailyData{
			TotalReports: totalSize,
			Page:         page,
			PageSize:     pageSize,
			DailyReports: dailyReports[offset:limit:limit],
		}, nil
	}
}

func (env *Env) ReadDailyReports(date string) (*DailyReport, error) {
	return env.ds.ReadDailyReport(date)
}
