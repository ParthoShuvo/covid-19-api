package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const configFilePath = "covid_19_api.json"

// Config holds configuration data.
type Config struct {
	configData *configData
	appName    string
}

// configData defines the COVID-19-API configuration
type configData struct {
	Description string
	Name        string
	AllowCORS   bool
	CORS        *CORSDef `json:"CORS"`
	Server      *ServerDef
	Logging     *LogDef
	Db          *DbDef `json:"Dataset"`
}

// CORSDef defines allowed cros settings
type CORSDef struct {
	AllowedOrigins   []string `json:"AllowedOrigins"`
	AllowCredentials bool     `json:"AllowCredentials"`
	AllowedMethods   []string `json:"AllowedMethods"`
	Debug            bool     `json:"Debug"`
}

// ServerDef defines a server address and port.
type ServerDef struct {
	Bind string
	Port int
}

// LogDef defines logging
type LogDef struct {
	Filename string
	Level    string
}

// DbDef database definition
type DbDef struct {
	CountryData   *CountryDataDef
	ArchievedData *DatasetDef
	CSSE          *DatasetDef
}

type CountryDataDef struct {
	CountryInfo    *DataDef
	CountryLatLong *DataDef
}

// DatasetDef dataset definition
type DatasetDef struct {
	DailyReports *DataDef
	TimeSeries   *TimeSeriesDataDef
}

// TimeSeriesDataDef defines confirmed, deaths, and recovered data definitions
type TimeSeriesDataDef struct {
	Confirmed *DataDef
	Deaths    *DataDef
	Recovered *DataDef
}

// DataDef defines a file data-type and path
type DataDef struct {
	Filetype string
	Filepath string
}

// NewConfig creates application configuration
func NewConfig(version string) *Config {
	cd, _ := loadConfig()
	an := cd.appName(version)
	return &Config{cd, an}
}

func loadConfig() (*configData, error) {
	jsonData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("failed to read file", configFilePath, err)
		return nil, err
	}

	var configData configData
	if err := json.Unmarshal(jsonData, &configData); err != nil {
		fmt.Println("failed to parse json")
		return nil, err
	}
	return &configData, nil
}

func (cd *configData) appName(version string) string {
	return fmt.Sprintf("%s/%s", cd.Name, version)
}

// AppName returns application name
func (c *Config) AppName() string {
	return c.appName
}

// AllowCORS determines whether cross origin calls are allowed.
func (c *Config) AllowCORS() bool {
	return c.configData.AllowCORS
}

// Database returns database definition
func (c *Config) Database() *DbDef {
	return c.configData.Db
}

//CROS definition
func (c *Config) CORS() *CORSDef {
	return c.configData.CORS
}

// Server returns the address and port to use for this service
func (c *Config) Server() *ServerDef {
	return c.configData.Server
}

func (s *ServerDef) String() string {
	return fmt.Sprintf("%s:%d", s.Bind, s.Port)
}

// Logging returns logfile and log level
func (c *Config) Logging() LogDef {
	return *c.configData.Logging
}
