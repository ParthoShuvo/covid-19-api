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
	Server      *serverDef
}

// serverDef defines a server address and port.
type serverDef struct {
	Bind string
	Port int
}

// Config holds configuration data.
type Config struct {
	configData *configData
	appName    string
}

func NewConfig(version string) *Config {
	cd, _ := loadConfig()
	an := cd.appName(version)
	return &Config{cd, an}
}

func (cd *configData) appName(version string) string {
	return fmt.Sprintf("%s/%s", cd.Name, version)
}

func (c *Config) AppName() string {
	return c.appName
}

func loadConfig() (*configData, error) {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("failed to read file", configFilePath, err)
		return nil, err
	}

	var configData configData
	if err := json.Unmarshal(data, &configData); err != nil {
		fmt.Println("failed to parse json")
		return nil, err
	}
	return &configData, nil
}
