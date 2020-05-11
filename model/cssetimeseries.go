package model

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/ParthoShuvo/fpingo/collection/list"
	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/cfg"
	"github.io/covid-19-api/errors"
	"github.io/covid-19-api/uc/cssetimeseries"
)

type (
	timeSeriesBuilder struct {
		timeSeriesMap cssetimeseries.TimeSeries
		dataSetter    func(d string, src *cssetimeseries.Report)
		err           error
		selectedTimes []string
	}

	csseTimeSeriesDate struct {
		date       string
		actualIndx int
	}

	timeSeriesCovidInfoValConsumer func(d string, src *cssetimeseries.Report)
)

func newTimeSeriesBuilder(selectedTimes []string) *timeSeriesBuilder {
	return &timeSeriesBuilder{
		timeSeriesMap: cssetimeseries.TimeSeries{},
		selectedTimes: selectedTimes,
	}
}

func (tb *timeSeriesBuilder) setConfirmedData(d string, src *cssetimeseries.Report) {
	val, err := strconv.Atoi(d)
	if err != nil {
		log.Error(err)
		return
	}
	src.Confirmed = val
}

func (tb *timeSeriesBuilder) setRecoveredData(d string, src *cssetimeseries.Report) {
	val, err := strconv.Atoi(d)
	if err != nil {
		log.Error(err)
		return
	}
	src.Recovered = val
}

func (tb *timeSeriesBuilder) setDeathsData(d string, src *cssetimeseries.Report) {
	val, err := strconv.Atoi(d)
	if err != nil {
		log.Error(err)
		return
	}
	src.Deaths = val
}

func (tb *timeSeriesBuilder) appendCovidInfoCases(filePath, caseName string, consumer timeSeriesCovidInfoValConsumer) *timeSeriesBuilder {
	tb.dataSetter = consumer
	timeSeriesMap, err := fetch(filePath, tb)
	if err != nil {
		return &timeSeriesBuilder{
			timeSeriesMap: tb.timeSeriesMap,
			err:           errors.InternalServerError.Wrapf(err, "failed to add %s cases", caseName),
			selectedTimes: tb.selectedTimes,
		}
	}
	return &timeSeriesBuilder{
		timeSeriesMap: timeSeriesMap.(cssetimeseries.TimeSeries),
		err:           tb.err,
		selectedTimes: tb.selectedTimes,
	}
}

func (tb *timeSeriesBuilder) parse(csvFile *os.File) (interface{}, error) {
	const (
		toLayout   string = "01-02-2006"
		fromLayout string = "1/2/06"
		dateCol    int    = 4
	)

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return tb.timeSeriesMap, errors.InternalServerError.Wrapf(err, "failed to parsed file: %s", csvFile.Name())
	}
	headerRow, dataRows := tb.split(records)
	dates := tb.parseDates(headerRow, dateCol, fromLayout, toLayout)

	for _, covidInfo := range dataRows {
		countryInfo := covidInfo[0:dateCol:dateCol]
		for _, d := range dates {
			reports, reportExists := tb.timeSeriesMap[cssetimeseries.Date(d.date)]
			if !reportExists {
				reports = []*cssetimeseries.Report{}
			}
			report, countryReportExists := tb.findCountryReport(reports, countryInfo)
			if !countryReportExists {
				report = tb.toReport(countryInfo)
			}
			tb.dataSetter(covidInfo[d.actualIndx], report)
			reports = append(reports, report)
			tb.timeSeriesMap[cssetimeseries.Date(d.date)] = reports
		}
	}
	return tb.timeSeriesMap, nil
}

func (tb *timeSeriesBuilder) split(records [][]string) (headerRow []string, dataRows [][]string) {
	headerRow, dataRows = records[0], records[1:]
	return
}

func (tb *timeSeriesBuilder) parseDates(data []string, startIndx int, fromlayout string, tolayout string) []*csseTimeSeriesDate {
	var dates []*csseTimeSeriesDate
	for i := startIndx; i < len(data); i++ {
		oldDate := data[i]
		od, err := parseDate(oldDate, fromlayout)
		if err != nil {
			log.Error(errors.NoType.Wrapf(err, "failed to parsed date: %s by layout: %s", oldDate, tolayout))
			continue
		}
		newDate := od.Format(tolayout)
		// fmt.Println(newDate)
		if !list.FromArray(tb.selectedTimes).Exists(func(t interface{}) bool { return t.(string) == newDate }) {
			continue
		}
		dates = append(dates, &csseTimeSeriesDate{newDate, i})
	}
	return dates
}

func (tb *timeSeriesBuilder) findOrNewReports(date string) []*cssetimeseries.Report {
	if reports, exists := tb.timeSeriesMap[cssetimeseries.Date(date)]; exists {
		return reports
	}
	var reports []*cssetimeseries.Report
	return reports
}

func (tb *timeSeriesBuilder) findCountryReport(reports []*cssetimeseries.Report, countryInfo []string) (*cssetimeseries.Report, bool) {
	const (
		province = iota
		country
		lat
		long
	)
	if report := list.FromArray(reports).FindFirst(func(r interface{}) bool {
		return r.(*cssetimeseries.Report).CountryRegion == countryInfo[country]
	}); report != nil {
		return report.(*cssetimeseries.Report), true
	}
	return &cssetimeseries.Report{}, false
}

func (tb *timeSeriesBuilder) toReport(countryInfo []string) *cssetimeseries.Report {
	const (
		province = iota
		country
		lat
		long
	)
	latitude, _ := strconv.ParseFloat(countryInfo[lat], 64)
	longitude, _ := strconv.ParseFloat(countryInfo[long], 64)
	return &cssetimeseries.Report{
		CountryRegion: countryInfo[country],
		ProvinceState: countryInfo[province],
		Lat:           latitude,
		Long:          longitude,
	}
}

func (tb *timeSeriesBuilder) build() (cssetimeseries.TimeSeries, error) {
	return tb.timeSeriesMap, tb.err
}

type csseTimeSeriesDataAccessor struct {
	dataDef *cfg.TimeSeriesDataDef
	builder *timeSeriesBuilder
}

func (da *csseTimeSeriesDataAccessor) GetAll() (cssetimeseries.TimeSeries, error) {
	timeSeries, err := da.combineAllCases()
	if err != nil {
		return cssetimeseries.TimeSeries{}, errors.InternalServerError.Wrap(err, "failed to get timeseries data")
	}
	return timeSeries, nil
}

func (da *csseTimeSeriesDataAccessor) combineAllCases() (cssetimeseries.TimeSeries, error) {
	confirmedCases := da.builder.appendCovidInfoCases(da.dataDef.Confirmed.Filepath, "Confirmed", da.builder.setConfirmedData)
	deathCases := confirmedCases.appendCovidInfoCases(da.dataDef.Deaths.Filepath, "Deaths", da.builder.setDeathsData)
	recoveredCases := deathCases.appendCovidInfoCases(da.dataDef.Recovered.Filepath, "Recovered", da.builder.setRecoveredData)
	return recoveredCases.build()
}

func (da *csseTimeSeriesDataAccessor) GetByTime(times []string) (cssetimeseries.TimeSeries, error) {
	timeSeriesMap, err := da.GetAll()
	if err != nil {
		return timeSeriesMap, err
	}
	selectedTimeSeriesMap := cssetimeseries.TimeSeries{}
	for _, time := range times {
		if reports, exists := timeSeriesMap[cssetimeseries.Date(time)]; exists {
			selectedTimeSeriesMap[cssetimeseries.Date(time)] = reports
		}
	}

	if len(selectedTimeSeriesMap) == 0 {
		return cssetimeseries.TimeSeries{}, errors.NotFound.Newf("selected timeseries ain't found in %v", times)
	}
	return selectedTimeSeriesMap, nil
}

func (database *DB) ReadCsseTimeSeries(times []string) (cssetimeseries.TimeSeries, error) {
	da := &csseTimeSeriesDataAccessor{database.config.CSSE.TimeSeries, newTimeSeriesBuilder(times)}
	return da.GetByTime(times)
}
