package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const configFilePath = "covid_19_api.json"

// configData defines the COVID-19-API configuration
type configData struct {
	Description string
	Name        string
	Server      *ServerDef
	Logging     *LogDef
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

// Config holds configuration data.
type Config struct {
	configData *configData
	appName    string
}

// NewConfig creates application configuration
func NewConfig(version string) *Config {
	cd, _ := loadConfig()
	an := cd.appName(version)
	return &Config{cd, an}
}

func (cd *configData) appName(version string) string {
	return fmt.Sprintf("%s/%s", cd.Name, version)
}

// AppName returns application name
func (c *Config) AppName() string {
	return c.appName
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
