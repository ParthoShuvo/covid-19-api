package resource

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/errors"
	"github.io/covid-19-api/resource/writer"
	"github.io/covid-19-api/uc/cssetimeseries"
)

type CsseTimeSeriesResource struct {
	env    *cssetimeseries.Env
	writer writer.Writer
}

const dateRequestLayout = "01-02-2006"

func NewCsseTimeSeriesResource(env *cssetimeseries.Env, writer writer.Writer) *CsseTimeSeriesResource {
	return &CsseTimeSeriesResource{env, writer}
}

func (res *CsseTimeSeriesResource) TimeSeriesFetcherInBetween() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		queries := req.URL.Query()
		parsedDates, parsedErr := res.parseDates(dateRequestLayout, queries.Get("start"), queries.Get("end"))
		if parsedErr != nil {
			log.Error(parsedErr)
			SendError(rw, parsedErr)
			return
		}
		startDate, endDate := parsedDates[0], parsedDates[1]
		if !startDate.Before(endDate) {
			err := errors.BadRequest.Newf("start date: %s is not before end date: %s", startDate, endDate)
			log.Error(err)
			SendError(rw, err)
			return
		}

		pinfo, pageInfoErr := parsePageInfo(queries)
		if pageInfoErr != nil {
			log.Error(pageInfoErr)
			SendError(rw, pageInfoErr)
			return
		}

		timeSeries, err := res.env.ReadTimeSeriesInBetweenWithPagination(startDate, endDate, pinfo.page, pinfo.pageSize)
		if err != nil {
			log.Error(err)
			SendError(rw, err)
			return
		}
		res.writer.Write(rw, timeSeries)
	}
}

func (res *CsseTimeSeriesResource) TimeSeriesFetcherByDate() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		dates := req.URL.Query()["date"]
		timeSeries, err := res.env.ReadTimeSeries(dates)
		if err != nil {
			log.Error(err)
			SendError(rw, err)
			return
		}
		res.writer.Write(rw, timeSeries)
	}
}

func (res *CsseTimeSeriesResource) parseDates(layout string, dates ...string) ([]time.Time, error) {
	var parsedDates []time.Time
	for _, d := range dates {
		parsedDate, err := res.parseDate(d, layout)
		if err != nil {
			return []time.Time{}, err
		}
		parsedDates = append(parsedDates, parsedDate)
	}
	return parsedDates, nil
}

func (res *CsseTimeSeriesResource) parseDate(date, layout string) (time.Time, error) {
	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, errors.BadRequest.Wrapf(err, "failed to parse date %s", date)
	}
	return parsedDate, nil
}
