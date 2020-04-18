package db

import (
	"encoding/csv"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.io/covid-19-api/cfg"
)

// MM-dd-YYYY.csv daily report fields
var dailyReportsField = struct {
	ProvinceState csvField
	CountryRegion csvField
	LastUpdate    csvField
	Confirmed     csvField
	Deaths        csvField
	Recovered     csvField
}{
	ProvinceState: 0,
	CountryRegion: 1,
	LastUpdate:    2,
	Confirmed:     3,
	Deaths:        4,
	Recovered:     5,
}

type (
	// CsseReport definition
	CsseReport struct {
		CountryRegion string `json:"country/region"`
		ProvinceState string `json:"province/state"`
		LastUpdate    string `json:"last-update"`
		Confirmed     int    `json:"confirmed"`
		Deaths        int    `json:"deaths"`
		Recovered     int    `json:"recovered"`
	}

	// CsseDailyReports definition
	CsseDailyReports struct {
		Date    string        `json:"-"`
		Reports []*CsseReport `json:"reports"`
	}
)

type csseDailyDataAccessor struct {
	dbConfig *cfg.DataDef
	parser   csseDailyReportsParser
}

func newCsseDailtDataAccessor() DataAccessor {
	return &csseDailyDataAccessor{dbConfig.CSSE.DailyReports, csseDailyReportsParser{}}
}

// GetAll returns all Countries
func (cda *csseDailyDataAccessor) GetAll() interface{} {
	filePaths := listFiles(cda.dbConfig.Filepath)
	return cda.getAllConcurrently(filePaths)
}

func (cda *csseDailyDataAccessor) getAllConcurrently(filePaths []string) map[string]*CsseDailyReports {
	records := make(chan *CsseDailyReports, len(filePaths))
	var waitgrp sync.WaitGroup
	waitgrp.Add(len(filePaths))
	for _, path := range filePaths {
		go func(path string, record chan<- *CsseDailyReports) {
			defer waitgrp.Done()
			result, _ := fetch(path, cda.parser)
			tmp := result.(*CsseDailyReports)
			record <- tmp
		}(path, records)
	}

	go func() {
		waitgrp.Wait()
		close(records)
	}()

	dailyReports := make(map[string]*CsseDailyReports)
	for record := range records {
		if record != nil {
			dailyReports[record.Date] = record
		}
	}
	return dailyReports
}

type csseDailyReportsParser struct{}

func (cdrp csseDailyReportsParser) parse(csvFile *os.File) (interface{}, error) {
	defer csvFile.Close()
	date := strings.TrimSuffix(filepath.Base(csvFile.Name()), path.Ext(csvFile.Name()))
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	reports := []*CsseReport{}
	for _, record := range records[1:] {
		deaths, _ := strconv.Atoi(record[dailyReportsField.Deaths])
		recovered, _ := strconv.Atoi(record[dailyReportsField.Recovered])
		confirmed, _ := strconv.Atoi(record[dailyReportsField.Confirmed])
		dailyReport := CsseReport{
			CountryRegion: record[dailyReportsField.CountryRegion],
			ProvinceState: record[dailyReportsField.ProvinceState],
			LastUpdate:    record[dailyReportsField.LastUpdate],
			Confirmed:     confirmed,
			Deaths:        deaths,
			Recovered:     recovered,
		}
		reports = append(reports, &dailyReport)
	}
	return &CsseDailyReports{date, reports}, nil
}
