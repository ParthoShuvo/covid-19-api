package db

import (
	"github.io/covid-19-api/cfg"
)

var dbConfig *cfg.DbDef

// AddConfig set database configuration
func AddConfig(dd *cfg.DbDef) {
	dbConfig = dd
}

type csvField int

// DataType definition
type DataType string

// Stored DataTypes
const (
	CountryData DataType = "CountryData"
)

// DataAccessor provies data accessing API or operations from high level business services
type DataAccessor interface {
	GetAll() interface{}
}

type newDataAccessor func() DataAccessor

// NewDataAccessor defines a dataAccessor according to DataType
func NewDataAccessor(dt DataType) DataAccessor {
	dataAccessorMap := map[DataType]newDataAccessor{
		CountryData: newCountryDataAccessor,
	}
	dataAccessor := dataAccessorMap[dt]
	return dataAccessor()
}
