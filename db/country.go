package db

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.io/covid-19-api/cfg"
)

// country_latlong.csv fields
var latlongField = struct {
	Country   csvField
	Latitude  csvField
	Longitude csvField
	Name      csvField
}{
	Country:   0,
	Latitude:  1,
	Longitude: 2,
	Name:      3,
}

// country_info.csv fields
var infoField = struct {
	Name                   csvField
	Alpha2                 csvField
	Alpha3                 csvField
	CountryCode            csvField
	ISO3166                csvField
	Region                 csvField
	SubRegion              csvField
	IntermediateRegion     csvField
	RegionCode             csvField
	SubRegionCode          csvField
	IntermediateRegionCode csvField
}{
	Name:                   0,
	Alpha2:                 1,
	Alpha3:                 2,
	CountryCode:            3,
	ISO3166:                4,
	Region:                 5,
	SubRegion:              6,
	IntermediateRegion:     7,
	RegionCode:             8,
	SubRegionCode:          9,
	IntermediateRegionCode: 10,
}

type (
	latlong struct {
		CC        string  `json:"-"`
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
		Name      string  `json:"-"`
	}

	info struct {
		Name                   string `json:"name"`
		Alpha2                 string `json:"cc"`
		Alpha3                 string `json:"-"`
		CountryCode            string `json:"-"`
		ISO3166                string `json:"-"`
		Region                 string `json:"region"`
		SubRegion              string `json:"sub-region"`
		IntermediateRegion     string `json:"-"`
		RegionCode             string `json:"-"`
		SubRegionCode          string `json:"-"`
		IntermediateRegionCode string `json:"-"`
	}

	// Country definition
	Country struct {
		Info    info
		LatLong latlong
	}
)

type countryLatLongParser struct{}

func (cllp countryLatLongParser) parse(csvFile *os.File) (interface{}, error) {
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	countryLatLongMap := make(map[string]latlong)
	for _, record := range records[1:] {
		latitude, _ := strconv.ParseFloat(record[latlongField.Latitude], 64)
		longitude, _ := strconv.ParseFloat(record[latlongField.Longitude], 64)
		countryLatLong := latlong{
			record[latlongField.Country],
			latitude,
			longitude,
			record[latlongField.Name],
		}
		countryLatLongMap[countryLatLong.CC] = countryLatLong
	}
	return countryLatLongMap, nil
}

type countryInfoParser struct{}

func (cip countryInfoParser) parse(csvFile *os.File) (interface{}, error) {
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	countryInfoMap := make(map[string]info)
	for _, record := range records[1:] {
		countryInfo := info{
			record[infoField.Name],
			record[infoField.Alpha2],
			record[infoField.Alpha3],
			record[infoField.CountryCode],
			record[infoField.ISO3166],
			record[infoField.Region],
			record[infoField.SubRegion],
			record[infoField.IntermediateRegion],
			record[infoField.RegionCode],
			record[infoField.SubRegionCode],
			record[infoField.IntermediateRegionCode],
		}
		countryInfoMap[countryInfo.Alpha2] = countryInfo
	}
	return countryInfoMap, nil
}

// CountryDataAccessor provies Countrydata accessing API and operations
type countryDataAccessor struct {
	dbConfig      *cfg.CountryDataDef
	infoParser    countryInfoParser
	latlongParser countryLatLongParser
}

// GetAll returns all Countries
func (cda *countryDataAccessor) GetAll() interface{} {
	countries := make([]Country, 0)
	i, err := fetch(cda.dbConfig.CountryInfo.Filepath, cda.infoParser)
	if err != nil {
		log.Fatal("Country Info parsing failed")
		return countries
	}
	infoMap := i.(map[string]info)

	i, err = fetch(cda.dbConfig.CountryLatLong.Filepath, cda.latlongParser)
	if err != nil {
		log.Fatal("Country LatLong parsing failed")
		return countries
	}
	latlongMap := i.(map[string]latlong)

	for cc, info := range infoMap {
		if latLong, exists := latlongMap[cc]; exists {
			country := Country{
				Info:    info,
				LatLong: latLong,
			}
			countries = append(countries, country)
		}
	}
	return countries
}

func newCountryDataAccessor() DataAccessor {
	return &countryDataAccessor{dbConfig.CountryData, countryInfoParser{}, countryLatLongParser{}}
}
