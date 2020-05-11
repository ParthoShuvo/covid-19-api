package cssetimeseries

import (
	"time"

	"github.io/covid-19-api/errors"
	"github.io/covid-19-api/uc"
)

const (
	layout          = "01-02-2006"
	defaultPage     = 1
	defaultPageSize = 10
)

type Env struct {
	ds CsseTimeSeriesDatastore
}

func New(ds CsseTimeSeriesDatastore) *Env {
	return &Env{ds}
}

func (env *Env) ReadTimeSeriesInBetween(startDate, endDate time.Time) (*CsseTimeSeriesData, error) {
	return env.ReadTimeSeriesInBetweenWithPagination(startDate, endDate, defaultPage, defaultPageSize)
}

func (env *Env) ReadTimeSeriesInBetweenWithPagination(startDate, endDate time.Time, page, pageSize int) (*CsseTimeSeriesData, error) {
	dates := env.generateDatesInBetween(startDate, endDate)
	timeSeries, err := env.ReadTimeSeries(dates)
	if err != nil {
		return &CsseTimeSeriesData{}, errors.NotFound.Newf("No time-series data found between %s and %s", startDate, endDate)
	}
	return env.paginateTimeSeries(timeSeries, page, pageSize, dates)
}

func (env *Env) ReadTimeSeries(times []string) (TimeSeries, error) {
	return env.ds.ReadCsseTimeSeries(times)
}

func (env *Env) generateDatesInBetween(startDate, endDate time.Time) []string {
	var dates []string
	date := startDate
	for {
		if date.After(endDate) {
			break
		}
		dates = append(dates, date.Format(layout))
		date = date.AddDate(0, 0, 1)
	}
	return dates
}

func (env *Env) paginateTimeSeries(timeSeries TimeSeries, page, pageSize int, dates []string) (*CsseTimeSeriesData, error) {
	paginator := uc.NewPaginator(defaultPageSize, defaultPage)
	totalSize := len(timeSeries)
	offset, limit := paginator.Paginate(totalSize, pageSize, page)
	paginatedTimeSeries := TimeSeries{}
	for item, date := offset, 0; item < limit; date++ {
		reports, exists := timeSeries[Date(dates[date])]
		if !exists {
			continue
		}
		item++
		paginatedTimeSeries[Date(dates[date])] = reports
	}
	return &CsseTimeSeriesData{
		Total:      totalSize,
		Page:       page,
		PageSize:   pageSize,
		TimeSeries: paginatedTimeSeries,
	}, nil
}
