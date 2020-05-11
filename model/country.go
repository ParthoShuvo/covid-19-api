package model

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/ParthoShuvo/covid-19-api/cfg"
	"github.com/ParthoShuvo/covid-19-api/errors"
	countryStore "github.com/ParthoShuvo/covid-19-api/uc/country"
	"github.com/ParthoShuvo/fpingo/collection/list"
	fn "github.com/ParthoShuvo/fpingo/util"
	"github.com/gocarina/gocsv"
)

type (
	Latlong struct {
		CC        string `csv:"country"`
		Latitude  string `csv:"latitude"`
		Longitude string `csv:"longitude"`
		Name      string `csv:"name"`
	}

	CountryInfo struct {
		Name                   string `csv:"name"`
		Alpha2                 string `csv:"alpha-2"`
		Alpha3                 string `csv:"alpha-3"`
		CountryCode            string `csv:"country-code"`
		ISO3166                string `csv:"iso_3166-2"`
		Region                 string `csv:"region"`
		SubRegion              string `csv:"sub-region"`
		IntermediateRegion     string `csv:"intermediate-region"`
		RegionCode             string `csv:"region-code"`
		SubRegionCode          string `csv:"sub-region-code"`
		IntermediateRegionCode string `csv:"intermediate-region-code"`
	}

	countryData struct {
		infoMap    map[ccKey]*CountryInfo
		latlongMap map[ccKey]*Latlong
	}
)

type ccKey string

type countryLatLongParser struct{}

func (cllp countryLatLongParser) parse(csvFile *os.File) (interface{}, error) {
	defer csvFile.Close()
	var (
		latlongs   []*Latlong
		latlongMap map[ccKey]*Latlong
	)
	if err := gocsv.UnmarshalCSV(csv.NewReader(csvFile), &latlongs); err != nil {
		return latlongMap, errors.InternalServerError.Wrapf(err, "failed to unmarshal %s csv file", csvFile.Name())
	}
	latlongMap = map[ccKey]*Latlong{}
	list.FromArray(latlongs).ForEach(func(i interface{}) {
		l := i.(*Latlong)
		latlongMap[ccKey(l.CC)] = l
	})
	return latlongMap, nil
}

type countryInfoParser struct{}

func (cip countryInfoParser) parse(csvFile *os.File) (interface{}, error) {
	defer csvFile.Close()
	var (
		infos   []*CountryInfo
		infoMap map[ccKey]*CountryInfo
	)
	if err := gocsv.UnmarshalCSV(csv.NewReader(csvFile), &infos); err != nil {
		return infoMap, errors.InternalServerError.Wrapf(err, "failed to unmarshal %s csv file", csvFile.Name())
	}
	infoMap = map[ccKey]*CountryInfo{}
	list.FromArray(infos).ForEach(func(i interface{}) {
		c := i.(*CountryInfo)
		infoMap[ccKey(c.Alpha2)] = c
	})
	return infoMap, nil
}

// countryDataAccessor provies Countrydata accessing API and operations
type countryDataAccessor struct {
	dbConfig      *cfg.CountryDataDef
	infoParser    countryInfoParser
	latlongParser countryLatLongParser
}

func (cda *countryDataAccessor) getAll() (*countryData, error) {
	i, err := fetch(cda.dbConfig.CountryInfo.Filepath, cda.infoParser)
	if err != nil {
		return &countryData{}, err
	}

	l, err := fetch(cda.dbConfig.CountryLatLong.Filepath, cda.latlongParser)
	if err != nil {
		return &countryData{}, err
	}
	infoMap, latlongMap := i.(map[ccKey]*CountryInfo), l.(map[ccKey]*Latlong)
	return &countryData{infoMap, latlongMap}, nil
}

func (database *DB) ReadAllCountries() ([]*countryStore.Country, error) {
	var countries []*countryStore.Country
	da := &countryDataAccessor{database.config.CountryData, countryInfoParser{}, countryLatLongParser{}}
	countryData, err := da.getAll()
	if err != nil {
		return countries, err
	}
	for cc, info := range countryData.infoMap {
		if latlong, exists := countryData.latlongMap[cc]; exists {
			countries = append(countries, database.toCountry(info, latlong))
		}
	}
	return countries, nil
}

func (db *DB) toCountry(info *CountryInfo, latlong *Latlong) *countryStore.Country {
	lat, _ := strconv.ParseFloat(latlong.Latitude, 64)
	lng, _ := strconv.ParseFloat(latlong.Longitude, 64)
	return &countryStore.Country{
		Name:                   info.Name,
		CC:                     latlong.CC,
		Alpha2:                 info.Alpha2,
		Alpha3:                 info.Alpha3,
		CountryCode:            info.CountryCode,
		ISO3166:                info.ISO3166,
		Region:                 info.Region,
		SubRegion:              info.SubRegion,
		IntermediateRegion:     info.IntermediateRegion,
		RegionCode:             info.RegionCode,
		SubRegionCode:          info.SubRegionCode,
		IntermediateRegionCode: info.IntermediateRegionCode,
		Latitude:               lat,
		Longitude:              lng,
	}
}

func (database *DB) ReadCountryByCC(cc string) (*countryStore.Country, error) {
	predicateFn := func(country interface{}) bool {
		return strings.EqualFold(country.(*countryStore.Country).CC, cc)
	}
	country, err := database.findCountry(predicateFn)
	if err != nil {
		err = errors.NotFound.Wrapf(err, "No country is found by CC: %s", cc)
		return &countryStore.Country{}, err
	}
	return country, nil
}

func (database *DB) ReadCountryByName(name string) (*countryStore.Country, error) {
	predicateFn := func(country interface{}) bool {
		return strings.EqualFold(country.(*countryStore.Country).Name, name)
	}
	country, err := database.findCountry(predicateFn)
	if err != nil {
		err = errors.NotFound.Wrapf(err, "No country is found by name: %s", name)
		return &countryStore.Country{}, err
	}
	return country, nil
}

func (database *DB) findCountry(p fn.Predicate) (*countryStore.Country, error) {
	countries, err := database.ReadAllCountries()
	if err != nil {
		return &countryStore.Country{}, err
	}
	i := list.FromArray(countries).FindFirst(p)
	if i == nil {
		return &countryStore.Country{}, errors.NotFound.New("No country found")
	}
	country := i.(*countryStore.Country)
	return country, nil
}
