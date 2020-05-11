package model

import (
	"encoding/csv"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ParthoShuvo/covid-19-api/cfg"
	"github.com/ParthoShuvo/covid-19-api/errors"
	"github.com/ParthoShuvo/covid-19-api/uc/country"
	"github.com/ParthoShuvo/covid-19-api/uc/cssedaily"
	"github.com/ParthoShuvo/fpingo/collection/list"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

type (
	// CsseRegionReport definition
	CsseRegionReport struct {
		ProvinceState string `csv:"Province/State"`
		CountryRegion string `csv:"Country/Region"`
		LastUpdate    string `csv:"Last Update"`
		Confirmed     string `csv:"Confirmed"`
		Deaths        string `csv:"Deaths"`
		Recovered     string `csv:"Recovered"`
	}

	csseDailyReport struct {
		date    string
		reports []*CsseRegionReport
	}

	csseDailyReportChan struct {
		dailyReport *csseDailyReport
		err         error
	}
)

type csseDailyReportsParser struct{}

func (parser csseDailyReportsParser) parse(csvFile *os.File) (interface{}, error) {
	defer csvFile.Close()
	var regionReports []*CsseRegionReport
	if err := gocsv.UnmarshalCSV(csv.NewReader(csvFile), &regionReports); err != nil {
		err = errors.InternalServerError.Wrapf(err, "failed to parse file: %s", csvFile.Name)
		return &csseDailyReport{}, err
	}
	return &csseDailyReport{
		date:    parser.parseDate(csvFile),
		reports: regionReports,
	}, nil
}

func (parser *csseDailyReportsParser) parseDate(file *os.File) string {
	return strings.TrimSuffix(filepath.Base(file.Name()), path.Ext(file.Name()))
}

type csseDailyDataAccessor struct {
	dbConfig *cfg.DataDef
	parser   csseDailyReportsParser
}

func (da *csseDailyDataAccessor) GetAll() ([]*csseDailyReport, error) {
	filePaths := listFiles(da.dbConfig.Filepath)
	return da.getAllConcurrently(filePaths)
}

func (da *csseDailyDataAccessor) getAllConcurrently(filePaths []string) ([]*csseDailyReport, error) {
	reportChan := make(chan *csseDailyReportChan, len(filePaths))
	da.concurrentRun(filePaths, reportChan)
	var (
		csseDailyReports []*csseDailyReport
		err              error
	)
	for report := range reportChan {
		if report.err != nil {
			err = errors.InternalServerError.Wrap(err, report.err.Error())
			continue
		}
		csseDailyReports = append(csseDailyReports, report.dailyReport)
	}
	if len(csseDailyReports) == 0 {
		return []*csseDailyReport{}, err
	}
	sort.Slice(csseDailyReports, da.prevDateComparator(csseDailyReports, "01-02-2006"))
	return csseDailyReports, nil
}

func (da *csseDailyDataAccessor) concurrentRun(filePaths []string, reportChan chan<- *csseDailyReportChan) {
	var waitgrp sync.WaitGroup
	waitgrp.Add(len(filePaths))
	for _, path := range filePaths {
		go func(path string) {
			defer waitgrp.Done()
			da.fetchAndSend(path, reportChan)
		}(path)
	}
	go func() {
		waitgrp.Wait()
		close(reportChan)
	}()
}

func (da *csseDailyDataAccessor) fetchAndSend(filePath string, sender chan<- *csseDailyReportChan) {
	if result, err := fetch(filePath, da.parser); err != nil {
		sender <- &csseDailyReportChan{&csseDailyReport{}, err}
	} else {
		sender <- &csseDailyReportChan{result.(*csseDailyReport), nil}
	}
}

func (da *csseDailyDataAccessor) prevDateComparator(dailyReports []*csseDailyReport, layout string) func(i, j int) bool {
	return func(i, j int) bool {
		prevDate, prevDateErr := parseDate(dailyReports[i].date, layout)
		nextDate, nextDateErr := parseDate(dailyReports[j].date, layout)
		if prevDateErr != nil {
			log.Error("failed to parsed date %s", dailyReports[i].date)
			return true
		}
		if nextDateErr != nil {
			log.Error("failed to parsed date %s", dailyReports[j].date)
			return true
		}
		return prevDate.Before(nextDate)
	}
}

func (database *DB) ReadAllDailyReports() ([]*cssedaily.DailyReport, error) {
	csseDa := &csseDailyDataAccessor{database.config.CSSE.DailyReports, csseDailyReportsParser{}}
	if csseDailyReports, err := csseDa.GetAll(); err != nil {
		err = errors.InternalServerError.Wrap(err, "failed to read csse daily reports")
		return []*cssedaily.DailyReport{}, err
	} else {
		countries, _ := database.ReadAllCountries()
		return database.composeCsseDailyWithCountries(csseDailyReports, countries), nil
	}
}

func (database *DB) ReadDailyReport(date string) (*cssedaily.DailyReport, error) {
	dailyReports, err := database.ReadAllDailyReports()
	if err != nil {
		return &cssedaily.DailyReport{}, err
	}
	report := list.FromArray(dailyReports).FindFirst(func(r interface{}) bool {
		return date == r.(*cssedaily.DailyReport).Date
	})
	if report == nil {
		return &cssedaily.DailyReport{}, errors.NotFound.Newf("csse daily report in not found on %s", date)
	}
	return report.(*cssedaily.DailyReport), nil
}

func (database *DB) composeCsseDailyWithCountries(csseDailyReports []*csseDailyReport, countries []*country.Country) []*cssedaily.DailyReport {
	countryList := list.FromArray(countries)
	regionMapper := func(d *CsseRegionReport) *cssedaily.RegionReport {
		region := database.toCsseRegionReport(d)
		c := countryList.FindFirst(func(c interface{}) bool {
			return strings.EqualFold(c.(*country.Country).Name, region.CountryRegion) ||
				strings.EqualFold(c.(*country.Country).CC, region.CountryRegion)
		})
		if c != nil {
			region.CC = c.(*country.Country).CC
		}
		return region
	}
	var dailyReports []*cssedaily.DailyReport
	for _, report := range csseDailyReports {
		dailyReports = append(dailyReports, database.toCsseDailyReport(report, regionMapper))
	}
	return dailyReports
}

func (database *DB) toCsseDailyReport(data *csseDailyReport, regionMapper func(data *CsseRegionReport) *cssedaily.RegionReport) *cssedaily.DailyReport {
	var (
		regionReports  []*cssedaily.RegionReport
		totalConfirmed int
		totalDeaths    int
		totalRecovered int
	)
	for _, d := range data.reports {
		regionReport := regionMapper(d)
		totalConfirmed += regionReport.Confirmed
		totalDeaths += regionReport.Deaths
		totalRecovered += regionReport.Recovered
		regionReports = append(regionReports, regionReport)
	}
	return &cssedaily.DailyReport{
		Date:           data.date,
		TotalConfirmed: totalConfirmed,
		TotalDeaths:    totalDeaths,
		TotalRecovered: totalDeaths,
		Reports:        regionReports,
	}
}

func (database *DB) toCsseRegionReport(data *CsseRegionReport) *cssedaily.RegionReport {
	confirmed, _ := strconv.Atoi(data.Confirmed)
	deaths, _ := strconv.Atoi(data.Deaths)
	recovered, _ := strconv.Atoi(data.Recovered)
	return &cssedaily.RegionReport{
		CountryRegion: data.CountryRegion,
		ProvinceState: data.ProvinceState,
		LastUpdate:    data.LastUpdate,
		Confirmed:     confirmed,
		Deaths:        deaths,
		Recovered:     recovered,
	}
}
