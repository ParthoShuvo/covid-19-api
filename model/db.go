package model

import (
	"github.io/covid-19-api/cfg"
	"github.io/covid-19-api/errors"
)

type DB struct {
	config *cfg.DbDef
}

func NewDB(dbConfig *cfg.DbDef) (*DB, error) {
	if dbConfig == nil || (cfg.DbDef{}) == *dbConfig {
		return &DB{}, errors.InternalServerError.New("DB config is empty")
	}
	return &DB{dbConfig}, nil
}
